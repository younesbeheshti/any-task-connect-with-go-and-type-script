package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	notifdomain "github.com/younesbeheshti/any-task-connect/backend/internal/notification/domain"
	paymentdomain "github.com/younesbeheshti/any-task-connect/backend/internal/payment/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/task/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/task/repository"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
)

const escrowFeePercent int64 = 8

// WalletService is the subset of wallet operations needed by tasks.
type WalletService interface {
	LockEscrow(ctx context.Context, input paymentdomain.EscrowLockInput) error
	ReleaseEscrow(ctx context.Context, input paymentdomain.EscrowReleaseInput) error
	RefundEscrow(ctx context.Context, taskID, requesterID uuid.UUID, amount, fee int64) error
}

// NotificationCreator creates in-app notifications. Injected post-construction
// to avoid an import cycle with the notification context.
type NotificationCreator interface {
	Create(ctx context.Context, input notifdomain.CreateNotificationInput) (*notifdomain.Notification, error)
}

// ApplicantSource returns the agents who applied to a task. Injected
// post-construction to avoid a cycle with the application context.
type ApplicantSource interface {
	AgentIDsForTask(ctx context.Context, taskID uuid.UUID) ([]uuid.UUID, error)
}

// TaskService implements service.Service.
type TaskService struct {
	repo       repository.Repository
	publisher  common.Publisher
	walletSvc  WalletService
	notifier   NotificationCreator
	applicants ApplicantSource
}

// NewTaskService creates a new TaskService.
func NewTaskService(repo repository.Repository, publisher common.Publisher) *TaskService {
	return &TaskService{repo: repo, publisher: publisher}
}

// SetWalletService injects the wallet service after construction to avoid import cycles.
func (s *TaskService) SetWalletService(ws WalletService) {
	s.walletSvc = ws
}

// SetNotifier injects the notification service after construction.
func (s *TaskService) SetNotifier(n NotificationCreator) {
	s.notifier = n
}

// SetApplicantSource injects the application service after construction.
func (s *TaskService) SetApplicantSource(a ApplicantSource) {
	s.applicants = a
}

func (s *TaskService) Create(ctx context.Context, input domain.CreateTaskInput) (*domain.Task, *domain.EscrowInfo, error) {
	publicID, err := s.repo.NextPublicID(ctx)
	if err != nil {
		return nil, nil, err
	}

	escrowFee := input.Budget * escrowFeePercent / 100

	task := &domain.Task{
		ID:             uuid.New(),
		PublicID:       publicID,
		Title:          input.Title,
		Description:    input.Description,
		CategoryID:     input.CategoryID,
		CityID:         input.CityID,
		Budget:         input.Budget,
		EscrowFee:      escrowFee,
		Currency:       "IRT",
		Status:         domain.TaskStatusCreated,
		Deadline:       input.Deadline,
		RequesterID:    input.RequesterID,
		AttachmentURLs: input.AttachmentURLs,
	}

	if err := s.repo.Create(ctx, task); err != nil {
		return nil, nil, err
	}

	if s.walletSvc != nil {
		lockInput := paymentdomain.EscrowLockInput{
			RequesterID: input.RequesterID,
			TaskID:      task.ID,
			Budget:      input.Budget,
			Fee:         escrowFee,
		}
		if lockErr := s.walletSvc.LockEscrow(ctx, lockInput); lockErr != nil {
			_ = s.repo.SoftDelete(ctx, task.ID)
			return nil, nil, apperrors.New("INSUFFICIENT_FUNDS", "موجودی کیف پول کافی نیست", 422, apperrors.ErrValidation)
		}
	}

	// Transition to OPEN only after escrow is confirmed locked.
	task, err = s.Transition(ctx, task.ID, domain.TaskStatusOpen, input.RequesterID, "")
	if err != nil {
		return nil, nil, err
	}

	if s.publisher != nil {
		_ = s.publisher.Publish(ctx, common.EventTaskCreated, map[string]any{
			"taskId":      task.ID.String(),
			"publicId":    task.PublicID,
			"requesterId": task.RequesterID.String(),
		})
	}

	escrow := &domain.EscrowInfo{
		Fee:      escrowFee,
		Held:     input.Budget + escrowFee,
		Currency: "IRT",
	}

	return task, escrow, nil
}

func (s *TaskService) GetByID(ctx context.Context, id uuid.UUID, viewerID *uuid.UUID) (*domain.Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, common.ErrNotFound) {
		return nil, apperrors.New("NOT_FOUND", "تسک یافت نشد", 404, apperrors.ErrNotFound)
	}
	return task, err
}

func (s *TaskService) GetByPublicID(ctx context.Context, publicID string, viewerID *uuid.UUID) (*domain.Task, error) {
	task, err := s.repo.GetByPublicID(ctx, publicID)
	if errors.Is(err, common.ErrNotFound) {
		return nil, apperrors.New("NOT_FOUND", "تسک یافت نشد", 404, apperrors.ErrNotFound)
	}
	return task, err
}

func (s *TaskService) Update(ctx context.Context, taskID, requesterID uuid.UUID, input domain.UpdateTaskInput) (*domain.Task, error) {
	task, err := s.repo.GetByID(ctx, taskID)
	if errors.Is(err, common.ErrNotFound) {
		return nil, apperrors.New("NOT_FOUND", "تسک یافت نشد", 404, apperrors.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	if task.RequesterID != requesterID {
		return nil, apperrors.New("FORBIDDEN", "دسترسی مجاز نیست", 403, apperrors.ErrForbidden)
	}
	if task.Status != domain.TaskStatusCreated && task.Status != domain.TaskStatusOpen {
		return nil, apperrors.New("INVALID_TRANSITION", "ویرایش تسک در این مرحله امکان‌پذیر نیست", 422, apperrors.ErrConflict)
	}

	if input.Title != nil {
		task.Title = *input.Title
	}
	if input.Description != nil {
		task.Description = *input.Description
	}
	if input.CategoryID != nil {
		task.CategoryID = *input.CategoryID
	}
	if input.CityID != nil {
		task.CityID = *input.CityID
	}
	if input.Budget != nil {
		task.Budget = *input.Budget
		task.EscrowFee = task.Budget * escrowFeePercent / 100
	}
	if input.Deadline != nil {
		task.Deadline = *input.Deadline
	}
	if input.AttachmentURLs != nil {
		task.AttachmentURLs = input.AttachmentURLs
	}

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) Delete(ctx context.Context, taskID, requesterID uuid.UUID) error {
	task, err := s.repo.GetByID(ctx, taskID)
	if errors.Is(err, common.ErrNotFound) {
		return apperrors.New("NOT_FOUND", "تسک یافت نشد", 404, apperrors.ErrNotFound)
	}
	if err != nil {
		return err
	}
	if task.RequesterID != requesterID {
		return apperrors.New("FORBIDDEN", "دسترسی مجاز نیست", 403, apperrors.ErrForbidden)
	}
	return s.repo.SoftDelete(ctx, taskID)
}

func (s *TaskService) Cancel(ctx context.Context, taskID, actorID uuid.UUID) (*domain.Task, error) {
	task, err := s.repo.GetByID(ctx, taskID)
	if errors.Is(err, common.ErrNotFound) {
		return nil, apperrors.New("NOT_FOUND", "تسک یافت نشد", 404, apperrors.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}

	updated, err := s.Transition(ctx, taskID, domain.TaskStatusCancelled, actorID, "لغو شده")
	if err != nil {
		return nil, err
	}

	if s.walletSvc != nil {
		_ = s.walletSvc.RefundEscrow(ctx, task.ID, task.RequesterID, task.Budget, task.EscrowFee)
	}

	// Let the involved agent(s) know the request was cancelled.
	s.notifyTaskCancelled(ctx, task)

	return updated, nil
}

// notifyTaskCancelled best-effort notifies the assigned agent, or — if the task
// was still open — every agent who had applied, that the task was cancelled.
func (s *TaskService) notifyTaskCancelled(ctx context.Context, task *domain.Task) {
	if s.notifier == nil {
		return
	}
	recipients := map[uuid.UUID]struct{}{}
	if task.AssignedAgentID != nil {
		recipients[*task.AssignedAgentID] = struct{}{}
	} else if s.applicants != nil {
		if ids, err := s.applicants.AgentIDsForTask(ctx, task.ID); err == nil {
			for _, id := range ids {
				recipients[id] = struct{}{}
			}
		}
	}
	if len(recipients) == 0 {
		return
	}
	deepLink := "/app/tasks/" + task.PublicID
	body := fmt.Sprintf("درخواست «%s» (%s) توسط درخواست‌دهنده لغو شد.", task.Title, task.PublicID)
	for agentID := range recipients {
		_, _ = s.notifier.Create(ctx, notifdomain.CreateNotificationInput{
			UserID:   agentID,
			Title:    "لغو درخواست",
			Body:     body,
			Type:     notifdomain.NotificationTypeSystem,
			DeepLink: &deepLink,
		})
	}
}

func (s *TaskService) Publish(ctx context.Context, taskID, requesterID uuid.UUID) (*domain.Task, error) {
	return s.Transition(ctx, taskID, domain.TaskStatusOpen, requesterID, "")
}

func (s *TaskService) Start(ctx context.Context, taskID, agentID uuid.UUID) (*domain.Task, error) {
	return s.Transition(ctx, taskID, domain.TaskStatusInProgress, agentID, "")
}

func (s *TaskService) Complete(ctx context.Context, taskID, agentID uuid.UUID) (*domain.Task, error) {
	task, err := s.Transition(ctx, taskID, domain.TaskStatusCompleted, agentID, "")
	if err != nil {
		return nil, err
	}
	return s.Transition(ctx, task.ID, domain.TaskStatusWaitingForVerification, agentID, "")
}

func (s *TaskService) Verify(ctx context.Context, taskID, requesterID uuid.UUID) (*domain.Task, error) {
	task, err := s.repo.GetByID(ctx, taskID)
	if errors.Is(err, common.ErrNotFound) {
		return nil, apperrors.New("NOT_FOUND", "تسک یافت نشد", 404, apperrors.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}

	// Requester verifies the work → release escrow to the agent. The task stays
	// at VERIFIED until the agent acknowledges receipt (ConfirmPayment → PAID).
	updated, transErr := s.Transition(ctx, taskID, domain.TaskStatusVerified, requesterID, "")
	if transErr != nil {
		return nil, transErr
	}

	if s.walletSvc != nil && task.AssignedAgentID != nil {
		releaseInput := paymentdomain.EscrowReleaseInput{
			TaskID:      task.ID,
			RequesterID: task.RequesterID,
			AgentID:     *task.AssignedAgentID,
			Amount:      task.Budget,
			Fee:         task.EscrowFee,
		}
		_ = s.walletSvc.ReleaseEscrow(ctx, releaseInput)
	}

	// Tell the agent the payment was released so they can confirm receipt.
	if s.notifier != nil && task.AssignedAgentID != nil {
		deepLink := "/app/tasks/" + task.PublicID
		_, _ = s.notifier.Create(ctx, notifdomain.CreateNotificationInput{
			UserID:   *task.AssignedAgentID,
			Title:    "پرداخت آزاد شد",
			Body:     fmt.Sprintf("درخواست «%s» (%s) تایید و مبلغ آزاد شد. لطفاً دریافت وجه را تایید کنید.", task.Title, task.PublicID),
			Type:     notifdomain.NotificationTypePayment,
			DeepLink: &deepLink,
		})
	}

	return updated, nil
}

// ConfirmPayment lets the assigned agent acknowledge they received the released
// payment, moving the task from VERIFIED to its final PAID state.
func (s *TaskService) ConfirmPayment(ctx context.Context, taskID, agentID uuid.UUID) (*domain.Task, error) {
	task, err := s.repo.GetByID(ctx, taskID)
	if errors.Is(err, common.ErrNotFound) {
		return nil, apperrors.New("NOT_FOUND", "تسک یافت نشد", 404, apperrors.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	if task.AssignedAgentID == nil || *task.AssignedAgentID != agentID {
		return nil, apperrors.New("FORBIDDEN", "دسترسی مجاز نیست", 403, apperrors.ErrForbidden)
	}
	if task.Status != domain.TaskStatusVerified {
		return nil, apperrors.New("CONFLICT", "این درخواست در وضعیت تایید‌شده نیست", 409, apperrors.ErrConflict)
	}

	updated, err := s.Transition(ctx, taskID, domain.TaskStatusPaid, agentID, "دریافت وجه تایید شد")
	if err != nil {
		return nil, err
	}

	// Let the requester know the agent confirmed receipt.
	if s.notifier != nil {
		deepLink := "/app/tasks/" + task.PublicID
		_, _ = s.notifier.Create(ctx, notifdomain.CreateNotificationInput{
			UserID:   task.RequesterID,
			Title:    "دریافت وجه تایید شد",
			Body:     fmt.Sprintf("مجری دریافت وجه درخواست «%s» (%s) را تایید کرد.", task.Title, task.PublicID),
			Type:     notifdomain.NotificationTypePayment,
			DeepLink: &deepLink,
		})
	}

	return updated, nil
}

func (s *TaskService) List(ctx context.Context, filter domain.TaskFilter, pg common.PaginationParams) (*common.PaginatedResult[domain.Task], error) {
	tasks, total, err := s.repo.List(ctx, filter, pg)
	if err != nil {
		return nil, err
	}
	page := pg.Page
	if page < 1 {
		page = 1
	}
	pageSize := pg.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	return &common.PaginatedResult[domain.Task]{
		Items:    tasks,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

func (s *TaskService) GetTimeline(ctx context.Context, taskID uuid.UUID) ([]domain.TaskTimeline, error) {
	return s.repo.GetTimeline(ctx, taskID)
}

func (s *TaskService) Transition(ctx context.Context, taskID uuid.UUID, next domain.TaskStatus, actorID uuid.UUID, note string) (*domain.Task, error) {
	task, err := s.repo.GetByID(ctx, taskID)
	if errors.Is(err, common.ErrNotFound) {
		return nil, apperrors.New("NOT_FOUND", "تسک یافت نشد", 404, apperrors.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}

	if !task.CanTransitionTo(next) {
		return nil, apperrors.New("INVALID_TRANSITION",
			"انتقال وضعیت غیرمجاز است", 422, apperrors.ErrConflict)
	}

	prev := task.Status
	if err := s.repo.UpdateStatus(ctx, task.ID, next); err != nil {
		return nil, err
	}
	task.Status = next
	task.UpdatedAt = time.Now()

	entry := &domain.TaskTimeline{
		ID:        uuid.New(),
		TaskID:    task.ID,
		ToStatus:  next,
		ActorID:   &actorID,
		Note:      note,
		CreatedAt: time.Now(),
	}
	if prev != "" {
		entry.FromStatus = &prev
	}
	_ = s.repo.AppendTimeline(ctx, entry)

	return task, nil
}

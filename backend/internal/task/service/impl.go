package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/task/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/task/repository"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
)

const escrowFeePercent int64 = 8

// TaskService implements service.Service.
type TaskService struct {
	repo      repository.Repository
	publisher common.Publisher
}

// NewTaskService creates a new TaskService.
func NewTaskService(repo repository.Repository, publisher common.Publisher) *TaskService {
	return &TaskService{repo: repo, publisher: publisher}
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

	// Immediately publish task to OPEN state.
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
	return s.Transition(ctx, taskID, domain.TaskStatusCancelled, actorID, "لغو شده")
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
	task, err := s.Transition(ctx, taskID, domain.TaskStatusVerified, requesterID, "")
	if err != nil {
		return nil, err
	}
	return s.Transition(ctx, task.ID, domain.TaskStatusPaid, requesterID, "")
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

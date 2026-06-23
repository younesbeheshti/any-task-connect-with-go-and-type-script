package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/application/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/application/repository"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	taskdomain "github.com/younesbeheshti/any-task-connect/backend/internal/task/domain"
	taskrepo "github.com/younesbeheshti/any-task-connect/backend/internal/task/repository"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
)

// ApplicationService implements service.Service.
type ApplicationService struct {
	repo     repository.Repository
	taskRepo taskrepo.Repository
}

// NewApplicationService creates a new ApplicationService.
func NewApplicationService(repo repository.Repository, taskRepo taskrepo.Repository) *ApplicationService {
	return &ApplicationService{repo: repo, taskRepo: taskRepo}
}

func (s *ApplicationService) Submit(ctx context.Context, input domain.SubmitApplicationInput) (*domain.Application, error) {
	// Check task exists and is open.
	task, err := s.taskRepo.GetByID(ctx, input.TaskID)
	if errors.Is(err, common.ErrNotFound) {
		return nil, apperrors.New("NOT_FOUND", "تسک یافت نشد", 404, apperrors.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	if task.Status != taskdomain.TaskStatusOpen {
		return nil, apperrors.New("CONFLICT", "امکان ثبت پیشنهاد برای این تسک وجود ندارد", 409, apperrors.ErrConflict)
	}

	// Check duplicate.
	existing, err := s.repo.GetByTaskAndAgent(ctx, input.TaskID, input.AgentID)
	if err == nil && existing != nil {
		return nil, apperrors.New("CONFLICT", "قبلاً برای این تسک پیشنهاد ثبت کرده‌اید", 409, apperrors.ErrConflict)
	}

	app := &domain.Application{
		ID:                     uuid.New(),
		TaskID:                 input.TaskID,
		AgentID:                input.AgentID,
		ProposalMessage:        input.ProposalMessage,
		ExpectedCompletionTime: input.ExpectedCompletionTime,
		ProposedPrice:          input.ProposedPrice,
		ETA:                    input.ETA,
		Status:                 domain.ApplicationStatusPending,
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}
	if err := s.repo.Create(ctx, app); err != nil {
		return nil, err
	}
	return app, nil
}

func (s *ApplicationService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Application, error) {
	app, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, common.ErrNotFound) {
		return nil, apperrors.New("NOT_FOUND", "درخواست یافت نشد", 404, apperrors.ErrNotFound)
	}
	return app, err
}

func (s *ApplicationService) ListByTask(ctx context.Context, taskID, requesterID uuid.UUID) ([]domain.Application, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if errors.Is(err, common.ErrNotFound) {
		return nil, apperrors.New("NOT_FOUND", "تسک یافت نشد", 404, apperrors.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	if task.RequesterID != requesterID {
		return nil, apperrors.New("FORBIDDEN", "دسترسی مجاز نیست", 403, apperrors.ErrForbidden)
	}
	return s.repo.ListByTask(ctx, taskID)
}

func (s *ApplicationService) ListByAgent(ctx context.Context, agentID uuid.UUID) ([]domain.Application, error) {
	return s.repo.ListByAgent(ctx, agentID)
}

func (s *ApplicationService) Accept(ctx context.Context, applicationID, requesterID uuid.UUID) (*domain.Application, error) {
	app, err := s.repo.GetByID(ctx, applicationID)
	if errors.Is(err, common.ErrNotFound) {
		return nil, apperrors.New("NOT_FOUND", "درخواست یافت نشد", 404, apperrors.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	if app.Status != domain.ApplicationStatusPending {
		return nil, apperrors.New("CONFLICT", "این درخواست قابل قبول کردن نیست", 409, apperrors.ErrConflict)
	}

	// Verify requester owns the task.
	task, err := s.taskRepo.GetByID(ctx, app.TaskID)
	if errors.Is(err, common.ErrNotFound) {
		return nil, apperrors.New("NOT_FOUND", "تسک یافت نشد", 404, apperrors.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	if task.RequesterID != requesterID {
		return nil, apperrors.New("FORBIDDEN", "دسترسی مجاز نیست", 403, apperrors.ErrForbidden)
	}

	// Accept this application.
	app.Status = domain.ApplicationStatusAccepted
	if err := s.repo.Update(ctx, app); err != nil {
		return nil, err
	}

	// Reject all other pending applications for this task.
	_ = s.repo.UpdateStatusByTask(ctx, app.TaskID, app.ID, domain.ApplicationStatusRejected)

	// Transition task to ASSIGNED.
	if err := s.taskRepo.UpdateStatus(ctx, task.ID, taskdomain.TaskStatusAssigned); err != nil {
		return nil, err
	}

	return app, nil
}

func (s *ApplicationService) Reject(ctx context.Context, applicationID, requesterID uuid.UUID) (*domain.Application, error) {
	app, err := s.repo.GetByID(ctx, applicationID)
	if errors.Is(err, common.ErrNotFound) {
		return nil, apperrors.New("NOT_FOUND", "درخواست یافت نشد", 404, apperrors.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	if app.Status != domain.ApplicationStatusPending {
		return nil, apperrors.New("CONFLICT", "این درخواست قابل رد کردن نیست", 409, apperrors.ErrConflict)
	}

	task, err := s.taskRepo.GetByID(ctx, app.TaskID)
	if errors.Is(err, common.ErrNotFound) {
		return nil, apperrors.New("NOT_FOUND", "تسک یافت نشد", 404, apperrors.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	if task.RequesterID != requesterID {
		return nil, apperrors.New("FORBIDDEN", "دسترسی مجاز نیست", 403, apperrors.ErrForbidden)
	}

	app.Status = domain.ApplicationStatusRejected
	if err := s.repo.Update(ctx, app); err != nil {
		return nil, err
	}
	return app, nil
}

func (s *ApplicationService) Withdraw(ctx context.Context, applicationID, agentID uuid.UUID) (*domain.Application, error) {
	app, err := s.repo.GetByID(ctx, applicationID)
	if errors.Is(err, common.ErrNotFound) {
		return nil, apperrors.New("NOT_FOUND", "درخواست یافت نشد", 404, apperrors.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	if app.AgentID != agentID {
		return nil, apperrors.New("FORBIDDEN", "دسترسی مجاز نیست", 403, apperrors.ErrForbidden)
	}
	if app.Status != domain.ApplicationStatusPending {
		return nil, apperrors.New("CONFLICT", "امکان پس‌گرفتن این درخواست وجود ندارد", 409, apperrors.ErrConflict)
	}

	app.Status = domain.ApplicationStatusWithdrawn
	if err := s.repo.Update(ctx, app); err != nil {
		return nil, err
	}
	return app, nil
}

package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/task/domain"
)

// Repository defines task persistence operations.
type Repository interface {
	Create(ctx context.Context, task *domain.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error)
	GetByPublicID(ctx context.Context, publicID string) (*domain.Task, error)
	Update(ctx context.Context, task *domain.Task) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.TaskStatus) error
	// AssignAgent sets the assigned agent and status atomically (used on application accept).
	AssignAgent(ctx context.Context, id, agentID uuid.UUID, status domain.TaskStatus) error
	// IncrementApplicantCount bumps the denormalized applicant_count by delta (e.g. +1 on apply).
	IncrementApplicantCount(ctx context.Context, id uuid.UUID, delta int) error
	List(ctx context.Context, filter domain.TaskFilter, pg common.PaginationParams) ([]domain.Task, int64, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error

	AppendTimeline(ctx context.Context, entry *domain.TaskTimeline) error
	GetTimeline(ctx context.Context, taskID uuid.UUID) ([]domain.TaskTimeline, error)

	NextPublicID(ctx context.Context) (string, error)
}

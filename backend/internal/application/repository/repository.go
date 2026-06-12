package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/application/domain"
)

// Repository defines application persistence operations.
type Repository interface {
	Create(ctx context.Context, app *domain.Application) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Application, error)
	GetByTaskAndAgent(ctx context.Context, taskID, agentID uuid.UUID) (*domain.Application, error)
	ListByTask(ctx context.Context, taskID uuid.UUID) ([]domain.Application, error)
	ListByAgent(ctx context.Context, agentID uuid.UUID) ([]domain.Application, error)
	Update(ctx context.Context, app *domain.Application) error
	UpdateStatusByTask(ctx context.Context, taskID uuid.UUID, exceptID uuid.UUID, status domain.ApplicationStatus) error
	CountByTask(ctx context.Context, taskID uuid.UUID) (int64, error)
}

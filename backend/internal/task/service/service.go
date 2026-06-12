package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/task/domain"
)

// Service defines task lifecycle use cases.
type Service interface {
	Create(ctx context.Context, input domain.CreateTaskInput) (*domain.Task, *domain.EscrowInfo, error)
	GetByID(ctx context.Context, id uuid.UUID, viewerID *uuid.UUID) (*domain.Task, error)
	GetByPublicID(ctx context.Context, publicID string, viewerID *uuid.UUID) (*domain.Task, error)
	Update(ctx context.Context, taskID, requesterID uuid.UUID, input domain.UpdateTaskInput) (*domain.Task, error)
	Cancel(ctx context.Context, taskID, actorID uuid.UUID) (*domain.Task, error)

	Publish(ctx context.Context, taskID, requesterID uuid.UUID) (*domain.Task, error)
	Start(ctx context.Context, taskID, agentID uuid.UUID) (*domain.Task, error)
	Complete(ctx context.Context, taskID, agentID uuid.UUID) (*domain.Task, error)
	Verify(ctx context.Context, taskID, requesterID uuid.UUID) (*domain.Task, error)

	List(ctx context.Context, filter domain.TaskFilter, pg common.PaginationParams) (*common.PaginatedResult[domain.Task], error)
	GetTimeline(ctx context.Context, taskID uuid.UUID) ([]domain.TaskTimeline, error)

	// Transition validates and applies a state change (internal use).
	Transition(ctx context.Context, taskID uuid.UUID, next domain.TaskStatus, actorID uuid.UUID, note string) (*domain.Task, error)
}

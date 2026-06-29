package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/application/domain"
)

// Service defines application use cases.
type Service interface {
	Submit(ctx context.Context, input domain.SubmitApplicationInput) (*domain.Application, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Application, error)
	ListByTask(ctx context.Context, taskPublicID string, requesterID uuid.UUID) ([]domain.Application, error)
	ListByAgent(ctx context.Context, agentID uuid.UUID) ([]domain.Application, error)
	Accept(ctx context.Context, applicationID, requesterID uuid.UUID) (*domain.Application, error)
	Reject(ctx context.Context, applicationID, requesterID uuid.UUID) (*domain.Application, error)
	Withdraw(ctx context.Context, applicationID, agentID uuid.UUID) (*domain.Application, error)
}

package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/rating/domain"
)

// Service defines review/rating use cases.
type Service interface {
	Create(ctx context.Context, input domain.CreateReviewInput) (*domain.Review, error)
	GetByTask(ctx context.Context, taskID uuid.UUID) ([]domain.Review, error)
	GetByUser(ctx context.Context, userID uuid.UUID) ([]domain.Review, error)
}

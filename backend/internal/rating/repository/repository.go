package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/rating/domain"
)

// Repository defines review persistence operations.
type Repository interface {
	Create(ctx context.Context, review *domain.Review) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Review, error)
	GetByTaskAndReviewer(ctx context.Context, taskID, reviewerID uuid.UUID) (*domain.Review, error)
	ListByReviewedUser(ctx context.Context, userID uuid.UUID) ([]domain.Review, error)
	ListByTask(ctx context.Context, taskID uuid.UUID) ([]domain.Review, error)
	AverageRating(ctx context.Context, userID uuid.UUID) (float64, int, error)
}

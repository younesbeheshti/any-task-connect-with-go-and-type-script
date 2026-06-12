package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/category/domain"
)

// Repository defines category persistence operations.
type Repository interface {
	Create(ctx context.Context, cat *domain.Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	GetByTitle(ctx context.Context, title string) (*domain.Category, error)
	List(ctx context.Context) ([]domain.Category, error)
	Update(ctx context.Context, cat *domain.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}

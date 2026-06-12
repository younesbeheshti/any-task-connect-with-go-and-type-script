package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/city/domain"
)

// Repository defines city persistence operations.
type Repository interface {
	Create(ctx context.Context, city *domain.City) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.City, error)
	GetByTitle(ctx context.Context, title string) (*domain.City, error)
	List(ctx context.Context) ([]domain.City, error)
	Update(ctx context.Context, city *domain.City) error
	Delete(ctx context.Context, id uuid.UUID) error
}

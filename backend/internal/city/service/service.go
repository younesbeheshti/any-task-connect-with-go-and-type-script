package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/city/domain"
)

// Service defines city use cases.
type Service interface {
	Create(ctx context.Context, title, province string) (*domain.City, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.City, error)
	List(ctx context.Context) ([]domain.City, error)
	Update(ctx context.Context, id uuid.UUID, title, province string) (*domain.City, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

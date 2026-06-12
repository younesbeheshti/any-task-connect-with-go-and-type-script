package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/category/domain"
)

// Service defines category use cases.
type Service interface {
	Create(ctx context.Context, title, icon string) (*domain.Category, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	List(ctx context.Context) ([]domain.Category, error)
	Update(ctx context.Context, id uuid.UUID, title, icon string) (*domain.Category, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

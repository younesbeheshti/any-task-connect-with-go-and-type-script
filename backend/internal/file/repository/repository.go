package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/file/domain"
)

// Repository persists file metadata.
type Repository interface {
	Create(ctx context.Context, f *domain.File) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.File, error)
}

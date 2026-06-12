package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/notification/domain"
)

// Repository defines notification persistence operations.
type Repository interface {
	Create(ctx context.Context, n *domain.Notification) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error)
	ListByUser(ctx context.Context, userID uuid.UUID, unreadOnly bool, pg common.PaginationParams) ([]domain.Notification, int64, error)
	MarkRead(ctx context.Context, id, userID uuid.UUID) error
	MarkAllRead(ctx context.Context, userID uuid.UUID) error
	CountUnread(ctx context.Context, userID uuid.UUID) (int64, error)
}

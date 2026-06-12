package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/notification/domain"
)

// Service defines notification use cases.
type Service interface {
	Create(ctx context.Context, input domain.CreateNotificationInput) (*domain.Notification, error)
	GetByID(ctx context.Context, id, userID uuid.UUID) (*domain.Notification, error)
	List(ctx context.Context, userID uuid.UUID, unreadOnly bool, pg common.PaginationParams) (*common.PaginatedResult[domain.Notification], error)
	MarkRead(ctx context.Context, id, userID uuid.UUID) error
	MarkAllRead(ctx context.Context, userID uuid.UUID) error
}

package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/withdraw/domain"
)

type Repository interface {
	Create(ctx context.Context, w *domain.WithdrawRequest) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.WithdrawRequest, error)
	ListByUser(ctx context.Context, userID uuid.UUID, pg common.PaginationParams) ([]domain.WithdrawRequest, int64, error)
	ListAll(ctx context.Context, pg common.PaginationParams) ([]domain.WithdrawRequest, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.Status, approvedBy *uuid.UUID) error
}

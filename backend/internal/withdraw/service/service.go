package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/withdraw/domain"
)

type Service interface {
	CreateWithdraw(ctx context.Context, input domain.CreateInput) (*domain.WithdrawRequest, error)
	GetByID(ctx context.Context, id, userID uuid.UUID) (*domain.WithdrawRequest, error)
	ListUserWithdraws(ctx context.Context, userID uuid.UUID, pg common.PaginationParams) (*common.PaginatedResult[domain.WithdrawRequest], error)
	Approve(ctx context.Context, id, adminID uuid.UUID) error
	Reject(ctx context.Context, id, adminID uuid.UUID) error
	ListAll(ctx context.Context, pg common.PaginationParams) (*common.PaginatedResult[domain.WithdrawRequest], error)
}

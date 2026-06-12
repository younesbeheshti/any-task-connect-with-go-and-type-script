package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/admin/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	paymentdomain "github.com/younesbeheshti/any-task-connect/backend/internal/payment/domain"
	userdomain "github.com/younesbeheshti/any-task-connect/backend/internal/user/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/user/repository"
)

// Service defines admin use cases.
type Service interface {
	GetMetrics(ctx context.Context) (*domain.DashboardMetrics, error)
	ListUsers(ctx context.Context, filter repository.UserFilter, pg common.PaginationParams) (*common.PaginatedResult[userdomain.User], error)
	UpdateUser(ctx context.Context, userID uuid.UUID, update domain.UserAdminUpdate) (*userdomain.User, error)
	ListTransactions(ctx context.Context, filter paymentdomain.TransactionFilter, pg common.PaginationParams) (*common.PaginatedResult[paymentdomain.Transaction], error)
	GetRevenueReport(ctx context.Context, rangeKey string) (*domain.RevenueReport, error)
	ResolveDispute(ctx context.Context, disputeID uuid.UUID, resolution domain.DisputeResolution) error
}

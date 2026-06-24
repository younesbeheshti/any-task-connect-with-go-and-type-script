package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/payment/domain"
	"gorm.io/gorm"
)

// Repository defines transaction ledger persistence.
type Repository interface {
	Create(ctx context.Context, tx *domain.Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error)
	GetByReference(ctx context.Context, ref string) (*domain.Transaction, error)
	List(ctx context.Context, filter domain.TransactionFilter, pg common.PaginationParams) ([]domain.Transaction, int64, error)
	ListAll(ctx context.Context, filter domain.TransactionFilter, pg common.PaginationParams) ([]domain.Transaction, int64, error)
	WithTx(tx *gorm.DB) Repository
}

package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/ledger/domain"
)

type Repository interface {
	Create(ctx context.Context, entries []domain.Entry) error
	GetByTransactionID(ctx context.Context, txID uuid.UUID) ([]domain.Entry, error)
}

package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/wallet/domain"
)

// Repository defines wallet persistence operations.
type Repository interface {
	Create(ctx context.Context, wallet *domain.Wallet) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Wallet, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Wallet, error)
	GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Wallet, error)
	Update(ctx context.Context, wallet *domain.Wallet) error
	LockFunds(ctx context.Context, walletID uuid.UUID, amount int64) error
	UnlockFunds(ctx context.Context, walletID uuid.UUID, amount int64) error
	ReleaseToAgent(ctx context.Context, requesterWalletID, agentWalletID uuid.UUID, amount int64) error
}

package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	paymentdomain "github.com/younesbeheshti/any-task-connect/backend/internal/payment/domain"
	walletdomain "github.com/younesbeheshti/any-task-connect/backend/internal/wallet/domain"
)

// Service defines wallet use cases.
type Service interface {
	EnsureWallet(ctx context.Context, userID uuid.UUID) (*walletdomain.Wallet, error)
	GetWallet(ctx context.Context, userID uuid.UUID) (*walletdomain.Wallet, error)
	// TopUp simulates a successful payment (mock gateway) and credits the wallet immediately.
	TopUp(ctx context.Context, userID uuid.UUID, amount int64) (*walletdomain.Wallet, *paymentdomain.Transaction, error)
	Withdraw(ctx context.Context, userID uuid.UUID, amount int64, iban string) error
	GetTransactions(ctx context.Context, userID uuid.UUID, filter paymentdomain.TransactionFilter, pg common.PaginationParams) (*common.PaginatedResult[paymentdomain.Transaction], error)
	GetStatistics(ctx context.Context, userID uuid.UUID) (*WalletStats, error)

	LockEscrow(ctx context.Context, input paymentdomain.EscrowLockInput) error
	ReleaseEscrow(ctx context.Context, input paymentdomain.EscrowReleaseInput) error
	RefundEscrow(ctx context.Context, taskID, requesterID uuid.UUID, amount, fee int64) error
}

type WalletStats struct {
	TotalEarned int64 `json:"totalEarned"`
	TotalSpent  int64 `json:"totalSpent"`
	TxCount     int64 `json:"txCount"`
}

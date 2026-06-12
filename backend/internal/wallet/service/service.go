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
	GetWallet(ctx context.Context, userID uuid.UUID) (*walletdomain.Wallet, error)
	TopUp(ctx context.Context, userID uuid.UUID, amount int64, cardID *uuid.UUID, returnURL string) (paymentURL string, err error)
	Withdraw(ctx context.Context, userID uuid.UUID, amount int64, iban string) error
	GetTransactions(ctx context.Context, userID uuid.UUID, filter paymentdomain.TransactionFilter, pg common.PaginationParams) (*common.PaginatedResult[paymentdomain.Transaction], error)
	EnsureWallet(ctx context.Context, userID uuid.UUID) (*walletdomain.Wallet, error)
}

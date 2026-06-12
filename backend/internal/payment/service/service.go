package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/payment/domain"
)

// Service defines escrow and payment use cases.
type Service interface {
	LockEscrow(ctx context.Context, input domain.EscrowLockInput) error
	ReleaseEscrow(ctx context.Context, input domain.EscrowReleaseInput) error
	RefundEscrow(ctx context.Context, taskID, requesterID uuid.UUID, amount int64) error
	ProcessTopUpCallback(ctx context.Context, reference string, success bool) error
}

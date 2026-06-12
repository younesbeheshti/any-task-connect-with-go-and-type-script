package domain

import (
	"time"

	"github.com/google/uuid"
)

// TransactionType categorizes ledger entries.
type TransactionType string

const (
	TransactionTypeEscrowLock    TransactionType = "ESCROW_LOCK"
	TransactionTypeEscrowRelease TransactionType = "ESCROW_RELEASE"
	TransactionTypeTopup         TransactionType = "TOPUP"
	TransactionTypeRefund        TransactionType = "REFUND"
	TransactionTypeWithdraw      TransactionType = "WITHDRAW"
	TransactionTypeFee           TransactionType = "FEE"
	TransactionTypeTransfer      TransactionType = "TRANSFER"
)

// TransactionStatus represents transaction processing state.
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "PENDING"
	TransactionStatusCompleted TransactionStatus = "COMPLETED"
	TransactionStatusFailed    TransactionStatus = "FAILED"
	TransactionStatusLocked    TransactionStatus = "LOCKED"
	TransactionStatusReleased  TransactionStatus = "RELEASED"
)

// Transaction is an immutable wallet ledger entry.
type Transaction struct {
	ID              uuid.UUID         `json:"id"`
	WalletID        uuid.UUID         `json:"-"`
	TaskID          *uuid.UUID        `json:"relatedTaskId"`
	Amount          int64             `json:"amount"`
	Type            TransactionType   `json:"type"`
	Status          TransactionStatus `json:"status"`
	ReferenceNumber string            `json:"id"`
	Description     string            `json:"description"`
	CreatedAt       time.Time         `json:"date"`
}

// TransactionFilter supports ledger queries.
type TransactionFilter struct {
	WalletID  *uuid.UUID
	TaskID    *uuid.UUID
	Type      *TransactionType
	DateFrom  *time.Time
	DateTo    *time.Time
}

// EscrowReleaseInput holds payment release parameters.
type EscrowReleaseInput struct {
	TaskID      uuid.UUID
	RequesterID uuid.UUID
	AgentID     uuid.UUID
	Amount      int64
	Fee         int64
}

// EscrowLockInput holds escrow lock parameters on task creation.
type EscrowLockInput struct {
	RequesterID uuid.UUID
	TaskID      uuid.UUID
	Budget      int64
	Fee         int64
}

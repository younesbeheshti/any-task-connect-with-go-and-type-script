package domain

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionTypeEscrowLock    TransactionType = "ESCROW_LOCK"
	TransactionTypeEscrowRelease TransactionType = "ESCROW_RELEASE"
	TransactionTypeTopup         TransactionType = "TOPUP"
	TransactionTypeDeposit       TransactionType = "DEPOSIT"
	TransactionTypeRefund        TransactionType = "REFUND"
	TransactionTypeWithdraw      TransactionType = "WITHDRAW"
	TransactionTypeFee           TransactionType = "FEE"
	TransactionTypeTransfer      TransactionType = "TRANSFER"
	TransactionTypePayment       TransactionType = "PAYMENT"
	TransactionTypeCommission    TransactionType = "COMMISSION"
	TransactionTypeAdjustment    TransactionType = "ADJUSTMENT"
)

func (t TransactionType) APIType() string {
	switch t {
	case TransactionTypeEscrowLock:
		return "escrow_in"
	case TransactionTypeEscrowRelease:
		return "release"
	case TransactionTypeTopup, TransactionTypeDeposit:
		return "topup"
	case TransactionTypeRefund:
		return "refund"
	case TransactionTypeWithdraw:
		return "withdraw"
	case TransactionTypeFee, TransactionTypeCommission:
		return "fee"
	default:
		return string(t)
	}
}

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "PENDING"
	TransactionStatusCompleted TransactionStatus = "COMPLETED"
	TransactionStatusFailed    TransactionStatus = "FAILED"
	TransactionStatusLocked    TransactionStatus = "LOCKED"
	TransactionStatusReleased  TransactionStatus = "RELEASED"
	TransactionStatusReversed  TransactionStatus = "REVERSED"
)

func (s TransactionStatus) APIStatus() string {
	switch s {
	case TransactionStatusPending:
		return "pending"
	case TransactionStatusCompleted:
		return "completed"
	case TransactionStatusFailed:
		return "failed"
	case TransactionStatusLocked:
		return "locked"
	case TransactionStatusReleased:
		return "released"
	case TransactionStatusReversed:
		return "reversed"
	default:
		return string(s)
	}
}

type Transaction struct {
	ID              uuid.UUID         `json:"id"`
	WalletID        uuid.UUID         `json:"-"`
	TaskID          *uuid.UUID        `json:"relatedTaskId"`
	Amount          int64             `json:"amount"`
	Currency        string            `json:"currency"`
	Type            TransactionType   `json:"type"`
	Status          TransactionStatus `json:"status"`
	ReferenceNumber string            `json:"referenceNumber"`
	Description     string            `json:"description"`
	CreatedAt       time.Time         `json:"date"`
}

type TransactionFilter struct {
	WalletID *uuid.UUID
	TaskID   *uuid.UUID
	Type     *TransactionType
	DateFrom *time.Time
	DateTo   *time.Time
}

type EscrowReleaseInput struct {
	TaskID      uuid.UUID
	RequesterID uuid.UUID
	AgentID     uuid.UUID
	Amount      int64
	Fee         int64
}

type EscrowLockInput struct {
	RequesterID uuid.UUID
	TaskID      uuid.UUID
	Budget      int64
	Fee         int64
}

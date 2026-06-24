package domain

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending  Status = "PENDING"
	StatusApproved Status = "APPROVED"
	StatusRejected Status = "REJECTED"
	StatusPaid     Status = "PAID"
)

func (s Status) APIStatus() string {
	switch s {
	case StatusPending:
		return "pending"
	case StatusApproved:
		return "approved"
	case StatusRejected:
		return "rejected"
	case StatusPaid:
		return "paid"
	default:
		return string(s)
	}
}

type WithdrawRequest struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	WalletID    uuid.UUID
	Amount      int64
	Status      Status
	BankAccount string
	ShebaNumber string
	Description string
	ApprovedBy  *uuid.UUID
	ApprovedAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateInput struct {
	UserID      uuid.UUID
	WalletID    uuid.UUID
	Amount      int64
	BankAccount string
	ShebaNumber string
	Description string
}

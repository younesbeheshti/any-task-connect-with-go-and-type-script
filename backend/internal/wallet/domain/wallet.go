package domain

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	ID                     uuid.UUID `json:"id"`
	UserID                 uuid.UUID `json:"userId"`
	AvailableBalance       int64     `json:"availableBalance"`
	LockedBalance          int64     `json:"lockedBalance"`
	PendingWithdrawBalance int64     `json:"pendingWithdrawBalance"`
	TotalEarned            int64     `json:"totalEarned"`
	TotalSpent             int64     `json:"totalSpent"`
	Currency               string    `json:"currency"`
	CreatedAt              time.Time `json:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
}

func (w *Wallet) CanLock(amount int64) bool {
	return amount > 0 && w.AvailableBalance >= amount
}

func (w *Wallet) CanUnlock(amount int64) bool {
	return amount > 0 && w.LockedBalance >= amount
}

func (w *Wallet) CanWithdraw(amount int64) bool {
	return amount > 0 && w.AvailableBalance >= amount
}

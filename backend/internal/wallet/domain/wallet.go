package domain

import (
	"time"

	"github.com/google/uuid"
)

// Wallet holds a user's financial balance.
type Wallet struct {
	ID               uuid.UUID `json:"id"`
	UserID           uuid.UUID `json:"userId"`
	AvailableBalance int64     `json:"availableBalance"`
	LockedBalance    int64     `json:"lockedBalance"`
	Currency         string    `json:"currency"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// LockFunds moves amount from available to locked.
func (w *Wallet) CanLock(amount int64) bool {
	return amount > 0 && w.AvailableBalance >= amount
}

// UnlockFunds moves amount from locked back to available (refund).
func (w *Wallet) CanUnlock(amount int64) bool {
	return amount > 0 && w.LockedBalance >= amount
}

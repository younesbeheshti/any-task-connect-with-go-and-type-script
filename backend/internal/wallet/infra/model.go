package infra

import (
	"time"

	"github.com/google/uuid"
)

type WalletModel struct {
	ID                     uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID                 uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	AvailableBalance       int64     `gorm:"not null;default:0"`
	LockedBalance          int64     `gorm:"not null;default:0"`
	PendingWithdrawBalance int64     `gorm:"not null;default:0"`
	TotalEarned            int64     `gorm:"not null;default:0"`
	TotalSpent             int64     `gorm:"not null;default:0"`
	Currency               string    `gorm:"size:10;not null;default:'IRT'"`
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

func (WalletModel) TableName() string { return "wallets" }

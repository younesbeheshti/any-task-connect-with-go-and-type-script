package infra

import (
	"time"

	"github.com/google/uuid"
)

type LedgerEntryModel struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TransactionID uuid.UUID `gorm:"type:uuid;not null;index"`
	Account       string    `gorm:"size:50;not null"`
	Debit         int64     `gorm:"not null;default:0"`
	Credit        int64     `gorm:"not null;default:0"`
	BalanceAfter  int64     `gorm:"not null;default:0"`
	CreatedAt     time.Time
}

func (LedgerEntryModel) TableName() string { return "ledger_entries" }

package infra

import (
	"time"

	"github.com/google/uuid"
)

type TransactionModel struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	WalletID        uuid.UUID  `gorm:"type:uuid;not null;index"`
	TaskID          *uuid.UUID `gorm:"type:uuid"`
	Amount          int64      `gorm:"not null"`
	Currency        string     `gorm:"size:10;not null;default:'IRT'"`
	Type            string     `gorm:"size:30;not null"`
	Status          string     `gorm:"size:20;not null;default:'PENDING'"`
	ReferenceNumber string     `gorm:"size:50;uniqueIndex;not null"`
	Description     string     `gorm:"type:text"`
	CreatedAt       time.Time
}

func (TransactionModel) TableName() string { return "transactions" }

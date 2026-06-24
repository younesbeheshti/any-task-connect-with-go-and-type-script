package infra

import (
	"time"

	"github.com/google/uuid"
)

type WithdrawModel struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null;index"`
	WalletID    uuid.UUID  `gorm:"type:uuid;not null"`
	Amount      int64      `gorm:"not null"`
	Status      string     `gorm:"type:withdraw_status;not null;default:'PENDING'"`
	BankAccount string     `gorm:"size:30"`
	ShebaNumber string     `gorm:"size:26"`
	Description string     `gorm:"type:text;not null;default:''"`
	ApprovedBy  *uuid.UUID `gorm:"type:uuid"`
	ApprovedAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (WithdrawModel) TableName() string { return "withdraw_requests" }

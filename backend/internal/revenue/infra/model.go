package infra

import (
	"time"

	"github.com/google/uuid"
)

type RevenueModel struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TaskID    *uuid.UUID `gorm:"type:uuid"`
	Source    string     `gorm:"size:50;not null;default:'commission'"`
	Amount    int64      `gorm:"not null"`
	CreatedAt time.Time
}

func (RevenueModel) TableName() string { return "revenue_records" }

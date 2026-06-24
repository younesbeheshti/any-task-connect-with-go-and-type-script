package infra

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AuditLogModel struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	EntityType  string          `gorm:"size:50;not null"`
	EntityID    *uuid.UUID      `gorm:"type:uuid"`
	Action      string          `gorm:"size:50;not null"`
	PerformedBy *uuid.UUID      `gorm:"type:uuid"`
	OldValue    json.RawMessage `gorm:"type:jsonb"`
	NewValue    json.RawMessage `gorm:"type:jsonb"`
	CreatedAt   time.Time
}

func (AuditLogModel) TableName() string { return "audit_logs" }

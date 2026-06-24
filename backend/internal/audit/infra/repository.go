package infra

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/audit/domain"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(ctx context.Context, log *domain.AuditLog) error {
	var oldVal, newVal json.RawMessage
	if log.OldValue != nil {
		b, _ := json.Marshal(log.OldValue)
		oldVal = b
	}
	if log.NewValue != nil {
		b, _ := json.Marshal(log.NewValue)
		newVal = b
	}
	m := AuditLogModel{
		ID:          log.ID,
		EntityType:  log.EntityType,
		EntityID:    log.EntityID,
		Action:      log.Action,
		PerformedBy: log.PerformedBy,
		OldValue:    oldVal,
		NewValue:    newVal,
	}
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	log.ID = m.ID
	log.CreatedAt = m.CreatedAt
	return nil
}

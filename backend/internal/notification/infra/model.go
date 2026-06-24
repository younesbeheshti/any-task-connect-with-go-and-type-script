package infra

import (
	"time"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/notification/domain"
)

type NotificationModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Title     string    `gorm:"size:255;not null"`
	Body      string    `gorm:"type:text;not null;default:''"`
	Type      string    `gorm:"size:50;not null"`
	IsRead    bool      `gorm:"not null;default:false"`
	DeepLink  *string   `gorm:"size:500"`
	CreatedAt time.Time
}

func (NotificationModel) TableName() string { return "notifications" }

func toDomain(m NotificationModel) *domain.Notification {
	return &domain.Notification{
		ID:        m.ID,
		UserID:    m.UserID,
		Title:     m.Title,
		Body:      m.Body,
		Type:      domain.NotificationType(m.Type),
		IsRead:    m.IsRead,
		DeepLink:  m.DeepLink,
		CreatedAt: m.CreatedAt,
	}
}

package infra

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/chat/domain"
	"gorm.io/gorm"
)

type ChatMessageModel struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TaskID     uuid.UUID `gorm:"type:uuid;not null;index"`
	SenderID   uuid.UUID `gorm:"type:uuid;not null"`
	ReceiverID uuid.UUID `gorm:"type:uuid;not null"`
	Message    string    `gorm:"type:text;not null;default:''"`
	Attachment []byte    `gorm:"type:jsonb"`
	Seen       bool      `gorm:"not null;default:false"`
	ReadAt     *time.Time
	CreatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (ChatMessageModel) TableName() string { return "chat_messages" }

func toMessageDomain(m ChatMessageModel) *domain.ChatMessage {
	msg := &domain.ChatMessage{
		ID:         m.ID,
		TaskID:     m.TaskID,
		SenderID:   m.SenderID,
		ReceiverID: m.ReceiverID,
		Message:    m.Message,
		Seen:       m.Seen,
		ReadAt:     m.ReadAt,
		CreatedAt:  m.CreatedAt,
	}
	if len(m.Attachment) > 0 && string(m.Attachment) != "null" {
		var att domain.Attachment
		if err := json.Unmarshal(m.Attachment, &att); err == nil {
			msg.Attachment = &att
		}
	}
	return msg
}

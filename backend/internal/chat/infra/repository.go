package infra

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/chat/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) CreateMessage(ctx context.Context, msg *domain.ChatMessage) error {
	m := ChatMessageModel{
		ID:         msg.ID,
		TaskID:     msg.TaskID,
		SenderID:   msg.SenderID,
		ReceiverID: msg.ReceiverID,
		Message:    msg.Message,
		Seen:       msg.Seen,
	}
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	if msg.Attachment != nil {
		raw, err := json.Marshal(msg.Attachment)
		if err == nil {
			m.Attachment = raw
		}
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	msg.ID = m.ID
	msg.CreatedAt = m.CreatedAt
	return nil
}

func (r *GormRepository) GetMessageByID(ctx context.Context, id uuid.UUID) (*domain.ChatMessage, error) {
	var m ChatMessageModel
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toMessageDomain(m), nil
}

func (r *GormRepository) ListMessages(ctx context.Context, taskID uuid.UUID, before *uuid.UUID, limit int) ([]domain.ChatMessage, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	q := r.db.WithContext(ctx).Where("task_id = ?", taskID).Order("created_at desc").Limit(limit)
	if before != nil {
		var ref ChatMessageModel
		if err := r.db.WithContext(ctx).First(&ref, "id = ?", *before).Error; err == nil {
			q = q.Where("created_at < ?", ref.CreatedAt)
		}
	}
	var models []ChatMessageModel
	if err := q.Find(&models).Error; err != nil {
		return nil, err
	}
	// reverse so oldest-first
	out := make([]domain.ChatMessage, len(models))
	for i, m := range models {
		out[len(models)-1-i] = *toMessageDomain(m)
	}
	return out, nil
}

func (r *GormRepository) MarkSeen(ctx context.Context, taskID, receiverID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&ChatMessageModel{}).
		Where("task_id = ? AND receiver_id = ? AND seen = false", taskID, receiverID).
		Updates(map[string]any{"seen": true, "read_at": now}).Error
}

func (r *GormRepository) CountUnread(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&ChatMessageModel{}).
		Where("receiver_id = ? AND seen = false", userID).
		Count(&count).Error
	return int(count), err
}

func (r *GormRepository) ListChatsForUser(ctx context.Context, userID uuid.UUID) ([]domain.ChatSummary, error) {
	// Fetch distinct task conversations the user is part of, with latest message.
	type row struct {
		TaskID      uuid.UUID `gorm:"column:task_id"`
		LastMessage string    `gorm:"column:last_message"`
		Unread      int       `gorm:"column:unread"`
		UpdatedAt   time.Time `gorm:"column:updated_at"`
		TaskTitle   string    `gorm:"column:task_title"`
	}
	var rows []row
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			cm.task_id,
			(SELECT message FROM chat_messages WHERE task_id = cm.task_id ORDER BY created_at DESC LIMIT 1) AS last_message,
			(SELECT COUNT(*) FROM chat_messages WHERE task_id = cm.task_id AND receiver_id = ? AND seen = false) AS unread,
			(SELECT created_at FROM chat_messages WHERE task_id = cm.task_id ORDER BY created_at DESC LIMIT 1) AS updated_at,
			COALESCE(t.title, '') AS task_title
		FROM chat_messages cm
		JOIN tasks t ON t.id = cm.task_id
		WHERE (cm.sender_id = ? OR cm.receiver_id = ?)
		  AND cm.deleted_at IS NULL
		GROUP BY cm.task_id, t.title
		ORDER BY updated_at DESC
		LIMIT 50
	`, userID, userID, userID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]domain.ChatSummary, len(rows))
	for i, rw := range rows {
		out[i] = domain.ChatSummary{
			ID:          rw.TaskID,
			TaskID:      rw.TaskID,
			Name:        rw.TaskTitle,
			LastMessage: rw.LastMessage,
			Unread:      rw.Unread,
			UpdatedAt:   rw.UpdatedAt,
		}
	}
	return out, nil
}

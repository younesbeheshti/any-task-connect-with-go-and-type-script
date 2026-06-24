package infra

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/notification/domain"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(ctx context.Context, n *domain.Notification) error {
	m := NotificationModel{
		ID:       n.ID,
		UserID:   n.UserID,
		Title:    n.Title,
		Body:     n.Body,
		Type:     string(n.Type),
		IsRead:   n.IsRead,
		DeepLink: n.DeepLink,
	}
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	n.ID = m.ID
	n.CreatedAt = m.CreatedAt
	return nil
}

func (r *GormRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error) {
	var m NotificationModel
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDomain(m), nil
}

func (r *GormRepository) ListByUser(ctx context.Context, userID uuid.UUID, unreadOnly bool, pg common.PaginationParams) ([]domain.Notification, int64, error) {
	q := r.db.WithContext(ctx).Model(&NotificationModel{}).Where("user_id = ?", userID)
	if unreadOnly {
		q = q.Where("is_read = false")
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	var models []NotificationModel
	if err := q.Order("created_at desc").Offset(pg.Offset()).Limit(pg.Limit()).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.Notification, len(models))
	for i, m := range models {
		out[i] = *toDomain(m)
	}
	return out, count, nil
}

func (r *GormRepository) MarkRead(ctx context.Context, id, userID uuid.UUID) error {
	result := r.db.WithContext(ctx).Model(&NotificationModel{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return common.ErrNotFound
	}
	return nil
}

func (r *GormRepository) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&NotificationModel{}).
		Where("user_id = ? AND is_read = false", userID).
		Update("is_read", true).Error
}

func (r *GormRepository) CountUnread(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&NotificationModel{}).
		Where("user_id = ? AND is_read = false", userID).
		Count(&count).Error
	return count, err
}

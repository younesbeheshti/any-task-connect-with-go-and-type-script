package infra

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/rating/domain"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(ctx context.Context, review *domain.Review) error {
	m := ReviewModel{
		ID:             review.ID,
		TaskID:         review.TaskID,
		ReviewerID:     review.ReviewerID,
		ReviewedUserID: review.ReviewedUserID,
		Rating:         review.Rating,
		Comment:        review.Comment,
	}
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	review.ID = m.ID
	review.CreatedAt = m.CreatedAt
	return nil
}

func (r *GormRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Review, error) {
	var m ReviewModel
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDomain(m), nil
}

func (r *GormRepository) GetByTaskAndReviewer(ctx context.Context, taskID, reviewerID uuid.UUID) (*domain.Review, error) {
	var m ReviewModel
	err := r.db.WithContext(ctx).First(&m, "task_id = ? AND reviewer_id = ?", taskID, reviewerID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDomain(m), nil
}

func (r *GormRepository) ListByReviewedUser(ctx context.Context, userID uuid.UUID) ([]domain.Review, error) {
	var models []ReviewModel
	if err := r.db.WithContext(ctx).Where("reviewed_user_id = ?", userID).Order("created_at desc").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Review, len(models))
	for i, m := range models {
		out[i] = *toDomain(m)
	}
	return out, nil
}

func (r *GormRepository) ListByTask(ctx context.Context, taskID uuid.UUID) ([]domain.Review, error) {
	var models []ReviewModel
	if err := r.db.WithContext(ctx).Where("task_id = ?", taskID).Order("created_at asc").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Review, len(models))
	for i, m := range models {
		out[i] = *toDomain(m)
	}
	return out, nil
}

func (r *GormRepository) AverageRating(ctx context.Context, userID uuid.UUID) (float64, int, error) {
	type result struct {
		Avg   float64
		Count int
	}
	var res result
	err := r.db.WithContext(ctx).Model(&ReviewModel{}).
		Select("COALESCE(AVG(rating), 0) as avg, COUNT(*) as count").
		Where("reviewed_user_id = ?", userID).
		Scan(&res).Error
	return res.Avg, res.Count, err
}

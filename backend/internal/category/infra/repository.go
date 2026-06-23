package infra

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/category/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"gorm.io/gorm"
)

// GormRepository implements category/repository.Repository using GORM.
type GormRepository struct {
	db *gorm.DB
}

// NewGormRepository creates a new category GormRepository.
func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) List(ctx context.Context) ([]domain.Category, error) {
	var models []CategoryModel
	if err := r.db.WithContext(ctx).Order("title asc").Find(&models).Error; err != nil {
		return nil, err
	}
	cats := make([]domain.Category, len(models))
	for i, m := range models {
		cats[i] = toDomain(m)
	}
	return cats, nil
}

func (r *GormRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	var m CategoryModel
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	c := toDomain(m)
	return &c, nil
}

func (r *GormRepository) GetByTitle(ctx context.Context, title string) (*domain.Category, error) {
	var m CategoryModel
	err := r.db.WithContext(ctx).Where("title = ?", title).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	c := toDomain(m)
	return &c, nil
}

func (r *GormRepository) Create(ctx context.Context, cat *domain.Category) error {
	m := CategoryModel{
		ID:        cat.ID,
		Title:     cat.Title,
		Icon:      cat.Icon,
		CreatedAt: cat.CreatedAt,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	cat.ID = m.ID
	cat.CreatedAt = m.CreatedAt
	return nil
}

func (r *GormRepository) Update(ctx context.Context, cat *domain.Category) error {
	return r.db.WithContext(ctx).Model(&CategoryModel{}).Where("id = ?", cat.ID).
		Updates(map[string]any{"title": cat.Title, "icon": cat.Icon}).Error
}

func (r *GormRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&CategoryModel{}, "id = ?", id).Error
}

func toDomain(m CategoryModel) domain.Category {
	return domain.Category{
		ID:        m.ID,
		Title:     m.Title,
		Icon:      m.Icon,
		CreatedAt: m.CreatedAt,
	}
}

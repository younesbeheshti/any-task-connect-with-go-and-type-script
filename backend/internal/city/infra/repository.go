package infra

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/city/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"gorm.io/gorm"
)

// GormRepository implements city/repository.Repository using GORM.
type GormRepository struct {
	db *gorm.DB
}

// NewGormRepository creates a new city GormRepository.
func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) List(ctx context.Context) ([]domain.City, error) {
	var models []CityModel
	if err := r.db.WithContext(ctx).Order("title asc").Find(&models).Error; err != nil {
		return nil, err
	}
	cities := make([]domain.City, len(models))
	for i, m := range models {
		cities[i] = toDomain(m)
	}
	return cities, nil
}

func (r *GormRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.City, error) {
	var m CityModel
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

func (r *GormRepository) GetByTitle(ctx context.Context, title string) (*domain.City, error) {
	var m CityModel
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

func (r *GormRepository) Create(ctx context.Context, city *domain.City) error {
	m := CityModel{
		ID:        city.ID,
		Title:     city.Title,
		Province:  city.Province,
		CreatedAt: city.CreatedAt,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	city.ID = m.ID
	city.CreatedAt = m.CreatedAt
	return nil
}

func (r *GormRepository) Update(ctx context.Context, city *domain.City) error {
	return r.db.WithContext(ctx).Model(&CityModel{}).Where("id = ?", city.ID).
		Updates(map[string]any{"title": city.Title, "province": city.Province}).Error
}

func (r *GormRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&CityModel{}, "id = ?", id).Error
}

func toDomain(m CityModel) domain.City {
	return domain.City{
		ID:        m.ID,
		Title:     m.Title,
		Province:  m.Province,
		CreatedAt: m.CreatedAt,
	}
}

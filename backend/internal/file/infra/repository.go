package infra

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/file/domain"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(ctx context.Context, f *domain.File) error {
	m := FileModel{
		ID:           f.ID,
		OriginalName: f.OriginalName,
		StoredName:   f.StoredName,
		Path:         f.Path,
		Size:         f.Size,
		MimeType:     f.MimeType,
		UploadedBy:   f.UploadedBy,
	}
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	f.ID = m.ID
	f.CreatedAt = m.CreatedAt
	return nil
}

func (r *GormRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.File, error) {
	var m FileModel
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toDomain(m), nil
}

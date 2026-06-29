package infra

import (
	"time"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/file/domain"
)

type FileModel struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OriginalName string    `gorm:"size:512;not null"`
	StoredName   string    `gorm:"size:512;uniqueIndex;not null"`
	Path         string    `gorm:"size:1024;not null"`
	Size         int64     `gorm:"not null"`
	MimeType     string    `gorm:"size:255;not null"`
	UploadedBy   uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedAt    time.Time
}

func (FileModel) TableName() string { return "files" }

func toDomain(m FileModel) *domain.File {
	return &domain.File{
		ID:           m.ID,
		OriginalName: m.OriginalName,
		StoredName:   m.StoredName,
		Path:         m.Path,
		Size:         m.Size,
		MimeType:     m.MimeType,
		UploadedBy:   m.UploadedBy,
		CreatedAt:    m.CreatedAt,
	}
}

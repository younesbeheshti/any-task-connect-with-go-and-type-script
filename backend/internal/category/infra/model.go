package infra

import (
	"time"

	"github.com/google/uuid"
)

// CategoryModel is the GORM persistence model for categories.
type CategoryModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title     string    `gorm:"uniqueIndex;size:100;not null"`
	Icon      string    `gorm:"size:50;not null;default:''"`
	CreatedAt time.Time
}

func (CategoryModel) TableName() string { return "categories" }

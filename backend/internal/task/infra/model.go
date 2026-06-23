package infra

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CategoryModel is a minimal category reference for task associations.
// It maps to the same "categories" table as category/infra.CategoryModel.
type CategoryModel struct {
	ID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Title string    `gorm:"column:title"`
}

func (CategoryModel) TableName() string { return "categories" }

// CityModel is a minimal city reference for task associations.
// It maps to the same "cities" table as city/infra.CityModel.
type CityModel struct {
	ID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Title string    `gorm:"column:title"`
}

func (CityModel) TableName() string { return "cities" }

// TaskModel is the GORM persistence model for tasks.
type TaskModel struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PublicID         string         `gorm:"uniqueIndex;size:20;not null"`
	Title            string         `gorm:"size:120;not null"`
	Description      string         `gorm:"type:text;not null"`
	CategoryID       uuid.UUID      `gorm:"type:uuid;not null"`
	CityID           uuid.UUID      `gorm:"type:uuid;not null"`
	Budget           int64          `gorm:"not null"`
	EscrowFee        int64          `gorm:"not null;default:0"`
	Status           string         `gorm:"size:40;not null;default:'CREATED'"`
	Deadline         time.Time      `gorm:"type:date;not null"`
	RequesterID      uuid.UUID      `gorm:"type:uuid;not null"`
	AssignedAgentID  *uuid.UUID     `gorm:"type:uuid"`
	AttachmentURLs   string         `gorm:"type:text;not null;default:'[]'"`
	ViewCount        int            `gorm:"not null;default:0"`
	ApplicationCount int            `gorm:"not null;default:0"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
	Category         CategoryModel  `gorm:"foreignKey:CategoryID;references:ID"`
	City             CityModel      `gorm:"foreignKey:CityID;references:ID"`
}

func (TaskModel) TableName() string { return "tasks" }

// TaskTimelineModel records task status transitions.
type TaskTimelineModel struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TaskID     uuid.UUID  `gorm:"type:uuid;not null;index"`
	FromStatus *string    `gorm:"size:40"`
	ToStatus   string     `gorm:"size:40;not null"`
	ActorID    *uuid.UUID `gorm:"type:uuid"`
	Note       string     `gorm:"type:text"`
	CreatedAt  time.Time
}

func (TaskTimelineModel) TableName() string { return "task_timelines" }

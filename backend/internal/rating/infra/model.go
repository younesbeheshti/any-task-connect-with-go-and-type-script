package infra

import (
	"time"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/rating/domain"
)

type ReviewModel struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TaskID         uuid.UUID `gorm:"type:uuid;not null;index"`
	ReviewerID     uuid.UUID `gorm:"type:uuid;not null"`
	ReviewedUserID uuid.UUID `gorm:"type:uuid;not null;index"`
	Rating         int       `gorm:"not null"`
	Comment        string    `gorm:"type:text;not null;default:''"`
	CreatedAt      time.Time
}

func (ReviewModel) TableName() string { return "reviews" }

func toDomain(m ReviewModel) *domain.Review {
	return &domain.Review{
		ID:             m.ID,
		TaskID:         m.TaskID,
		ReviewerID:     m.ReviewerID,
		ReviewedUserID: m.ReviewedUserID,
		Rating:         m.Rating,
		Comment:        m.Comment,
		CreatedAt:      m.CreatedAt,
	}
}

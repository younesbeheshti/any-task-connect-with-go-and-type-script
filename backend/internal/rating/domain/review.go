package domain

import (
	"time"

	"github.com/google/uuid"
)

// Review is a post-task rating between participants.
type Review struct {
	ID             uuid.UUID `json:"id"`
	TaskID         uuid.UUID `json:"taskId"`
	ReviewerID     uuid.UUID `json:"reviewerId"`
	ReviewedUserID uuid.UUID `json:"reviewedUserId"`
	Rating         int       `json:"rating"`
	Comment        string    `json:"comment"`
	CreatedAt      time.Time `json:"createdAt"`
}

// CreateReviewInput holds review creation data.
type CreateReviewInput struct {
	TaskID         uuid.UUID
	ReviewerID     uuid.UUID
	ReviewedUserID uuid.UUID
	Rating         int
	Comment        string
}

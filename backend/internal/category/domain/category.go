package domain

import (
	"time"

	"github.com/google/uuid"
)

// Category is a task classification.
type Category struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"createdAt"`
}

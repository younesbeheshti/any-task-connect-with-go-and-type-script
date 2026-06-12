package domain

import (
	"time"

	"github.com/google/uuid"
)

// City is a geographic location for tasks.
type City struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Province  string    `json:"province"`
	CreatedAt time.Time `json:"createdAt"`
}

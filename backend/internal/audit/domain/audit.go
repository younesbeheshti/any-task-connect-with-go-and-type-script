package domain

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID          uuid.UUID
	EntityType  string
	EntityID    *uuid.UUID
	Action      string
	PerformedBy *uuid.UUID
	OldValue    map[string]any
	NewValue    map[string]any
	CreatedAt   time.Time
}

type CreateInput struct {
	EntityType  string
	EntityID    *uuid.UUID
	Action      string
	PerformedBy *uuid.UUID
	OldValue    map[string]any
	NewValue    map[string]any
}

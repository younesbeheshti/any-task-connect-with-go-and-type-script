package domain

import (
	"time"

	"github.com/google/uuid"
)

// ApplicationStatus represents application lifecycle.
type ApplicationStatus string

const (
	ApplicationStatusPending   ApplicationStatus = "PENDING"
	ApplicationStatusAccepted  ApplicationStatus = "ACCEPTED"
	ApplicationStatusRejected  ApplicationStatus = "REJECTED"
	ApplicationStatusWithdrawn ApplicationStatus = "WITHDRAWN"
)

// Application is an agent's proposal for a task.
type Application struct {
	ID                     uuid.UUID         `json:"id"`
	TaskID                 uuid.UUID         `json:"taskId"`
	AgentID                uuid.UUID         `json:"-"`
	ProposalMessage        string            `json:"message"`
	ExpectedCompletionTime *time.Time        `json:"-"`
	ProposedPrice          *int64            `json:"price"`
	ETA                    string            `json:"eta"`
	Status                 ApplicationStatus `json:"status"`
	CreatedAt              time.Time         `json:"createdAt"`
	UpdatedAt              time.Time         `json:"updatedAt"`
}

// SubmitApplicationInput holds application creation data.
type SubmitApplicationInput struct {
	TaskID                 uuid.UUID
	AgentID                uuid.UUID
	ProposalMessage        string
	ExpectedCompletionTime *time.Time
	ProposedPrice          *int64
	ETA                    string
}

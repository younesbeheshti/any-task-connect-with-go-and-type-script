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
	ProposalMessage        string            `json:"proposalMessage"`
	ExpectedCompletionTime *time.Time        `json:"-"`
	ProposedPrice          *int64            `json:"proposedPrice"`
	ETA                    string            `json:"eta"`
	Status                 ApplicationStatus `json:"status"`
	CreatedAt              time.Time         `json:"createdAt"`
	UpdatedAt              time.Time         `json:"updatedAt"`
	// Agent is populated when the application is loaded with its agent (e.g. the
	// requester reviewing applicants); nil on bare returns like Submit/Accept.
	Agent *AgentSummary `json:"agent,omitempty"`
	// Task is populated when listing an agent's own applications, so the UI can
	// show which task each proposal belongs to.
	Task *TaskSummary `json:"task,omitempty"`
}

// TaskSummary is the minimal task context shown alongside an agent's application.
type TaskSummary struct {
	ID       string `json:"id"` // public id (e.g. "TB-7")
	Title    string `json:"title"`
	Budget   int64  `json:"budget"`
	City     string `json:"city"`
	Deadline string `json:"deadline"`
	Status   string `json:"status"` // API status value
}

// AgentSummary is the public agent profile shown to a requester picking an applicant.
type AgentSummary struct {
	ID             uuid.UUID `json:"id"`
	FullName       string    `json:"fullName"`
	Rating         float64   `json:"rating"`
	CompletedCount int       `json:"completedCount"`
	IsVerified     bool      `json:"isVerified"`
}

// SubmitApplicationInput holds application creation data.
type SubmitApplicationInput struct {
	// TaskPublicID is the public task identifier (e.g. "TB-7") from the URL; the
	// service resolves it to the internal task UUID. Matches the public-ID URL
	// convention used by the task action endpoints.
	TaskPublicID           string
	AgentID                uuid.UUID
	ProposalMessage        string
	ExpectedCompletionTime *time.Time
	ProposedPrice          *int64
	ETA                    string
}

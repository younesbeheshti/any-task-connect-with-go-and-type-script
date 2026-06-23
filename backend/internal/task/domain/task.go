package domain

import (
	"time"

	"github.com/google/uuid"
)

// TaskStatus represents the internal task lifecycle state.
type TaskStatus string

const (
	TaskStatusCreated                TaskStatus = "CREATED"
	TaskStatusOpen                   TaskStatus = "OPEN"
	TaskStatusAssigned               TaskStatus = "ASSIGNED"
	TaskStatusInProgress             TaskStatus = "IN_PROGRESS"
	TaskStatusCompleted              TaskStatus = "COMPLETED"
	TaskStatusWaitingForVerification TaskStatus = "WAITING_FOR_VERIFICATION"
	TaskStatusVerified               TaskStatus = "VERIFIED"
	TaskStatusPaid                   TaskStatus = "PAID"
	TaskStatusCancelled              TaskStatus = "CANCELLED"
)

// ValidTransitions defines allowed state machine edges.
var ValidTransitions = map[TaskStatus][]TaskStatus{
	TaskStatusCreated:                {TaskStatusOpen, TaskStatusCancelled},
	TaskStatusOpen:                   {TaskStatusAssigned, TaskStatusCancelled},
	TaskStatusAssigned:               {TaskStatusInProgress, TaskStatusCancelled},
	TaskStatusInProgress:             {TaskStatusCompleted, TaskStatusCancelled},
	TaskStatusCompleted:              {TaskStatusWaitingForVerification},
	TaskStatusWaitingForVerification: {TaskStatusVerified},
	TaskStatusVerified:               {TaskStatusPaid},
}

// APIStatus maps domain status to frontend contract values.
func (s TaskStatus) APIStatus() string {
	switch s {
	case TaskStatusCreated:
		return "posted"
	case TaskStatusOpen:
		return "awaiting_applicants"
	case TaskStatusAssigned:
		return "accepted"
	case TaskStatusInProgress:
		return "in_progress"
	case TaskStatusCompleted:
		return "completed"
	case TaskStatusWaitingForVerification, TaskStatusVerified:
		return "awaiting_verification"
	case TaskStatusPaid:
		return "paid"
	case TaskStatusCancelled:
		return "cancelled"
	default:
		return string(s)
	}
}

// Task is the task aggregate root.
type Task struct {
	ID              uuid.UUID    `json:"-"`
	PublicID        string       `json:"id"` // TB-1042 format for API
	Title           string       `json:"title"`
	Description     string       `json:"description"`
	CategoryID      uuid.UUID    `json:"-"`
	CityID          uuid.UUID    `json:"-"`
	Category        string       `json:"category"`
	City            string       `json:"city"`
	Budget          int64        `json:"budget"`
	EscrowFee       int64        `json:"-"`
	Currency        string       `json:"currency"`
	Status          TaskStatus   `json:"status"`
	Deadline        time.Time    `json:"deadline"`
	RequesterID     uuid.UUID    `json:"-"`
	AssignedAgentID *uuid.UUID   `json:"assignedAgentId"`
	AttachmentURLs  []string     `json:"-"`
	ApplicantsCount int          `json:"applicantsCount"`
	CreatedAt       time.Time    `json:"createdAt"`
	UpdatedAt       time.Time    `json:"updatedAt"`
	DeletedAt       *time.Time   `json:"-"`
}

// CanTransitionTo validates whether a state transition is allowed.
func (t *Task) CanTransitionTo(next TaskStatus) bool {
	allowed, ok := ValidTransitions[t.Status]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == next {
			return true
		}
	}
	return false
}

// TaskTimeline records a status change for audit and UI display.
type TaskTimeline struct {
	ID         uuid.UUID   `json:"id"`
	TaskID     uuid.UUID   `json:"taskId"`
	FromStatus *TaskStatus `json:"fromStatus"`
	ToStatus   TaskStatus  `json:"toStatus"`
	ActorID    *uuid.UUID  `json:"actorId"`
	Note       string      `json:"note"`
	CreatedAt  time.Time   `json:"createdAt"`
}

// CreateTaskInput holds data for task creation.
type CreateTaskInput struct {
	Title          string
	Description    string
	CategoryID     uuid.UUID
	CityID         uuid.UUID
	Budget         int64
	Deadline       time.Time
	AttachmentURLs []string
	RequesterID    uuid.UUID
}

// UpdateTaskInput holds mutable task fields (pre-assignment only).
type UpdateTaskInput struct {
	Title          *string
	Description    *string
	CategoryID     *uuid.UUID
	CityID         *uuid.UUID
	Budget         *int64
	Deadline       *time.Time
	AttachmentURLs []string
}

// TaskFilter supports list/search queries.
type TaskFilter struct {
	CityID      *uuid.UUID
	CityTitle   string
	CategoryID  *uuid.UUID
	CatTitle    string
	Status      *TaskStatus
	Query       string
	Sort        string // newest, budget, deadline
	MinBudget   *int64
	MaxBudget   *int64
	RequesterID *uuid.UUID
	AgentID     *uuid.UUID
}

// EscrowInfo returned on task creation.
type EscrowInfo struct {
	Fee      int64  `json:"fee"`
	Held     int64  `json:"held"`
	Currency string `json:"currency"`
}

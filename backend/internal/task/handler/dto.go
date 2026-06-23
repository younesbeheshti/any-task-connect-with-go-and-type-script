package handler

import (
	"time"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/task/domain"
)

// TaskResponse is the API response shape for a task.
type TaskResponse struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Category        string     `json:"category"`
	City            string     `json:"city"`
	Budget          int64      `json:"budget"`
	Currency        string     `json:"currency"`
	Status          string     `json:"status"`
	Deadline        string     `json:"deadline"`
	RequesterID     uuid.UUID  `json:"requesterId"`
	AssignedAgentID *uuid.UUID `json:"assignedAgentId"`
	ApplicantsCount int        `json:"applicantsCount"`
	AttachmentURLs  []string   `json:"attachmentUrls"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

// EscrowResponse wraps escrow info in the create task response.
type EscrowResponse struct {
	Fee      int64  `json:"fee"`
	Held     int64  `json:"held"`
	Currency string `json:"currency"`
}

// CreateTaskResponse is the body returned on task creation.
type CreateTaskResponse struct {
	Task   TaskResponse   `json:"task"`
	Escrow EscrowResponse `json:"escrow"`
}

func toTaskResponse(t *domain.Task) TaskResponse {
	urls := t.AttachmentURLs
	if urls == nil {
		urls = []string{}
	}
	return TaskResponse{
		ID:              t.PublicID,
		Title:           t.Title,
		Description:     t.Description,
		Category:        t.Category,
		City:            t.City,
		Budget:          t.Budget,
		Currency:        t.Currency,
		Status:          t.Status.APIStatus(),
		Deadline:        t.Deadline.Format("2006-01-02"),
		RequesterID:     t.RequesterID,
		AssignedAgentID: t.AssignedAgentID,
		ApplicantsCount: t.ApplicantsCount,
		AttachmentURLs:  urls,
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
	}
}

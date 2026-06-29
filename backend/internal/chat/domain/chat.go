package domain

import (
	"time"

	"github.com/google/uuid"
)

// ChatMessage represents a single message in a task conversation.
type ChatMessage struct {
	ID          uuid.UUID  `json:"id"`
	TaskID      uuid.UUID  `json:"-"`
	SenderID    uuid.UUID  `json:"from"`
	ReceiverID  uuid.UUID  `json:"-"`
	Message     string     `json:"text"`
	Attachment  *Attachment `json:"attachments"`
	Seen        bool       `json:"-"`
	ReadAt      *time.Time `json:"readAt"`
	CreatedAt   time.Time  `json:"createdAt"`
}

// Attachment holds file metadata for chat messages.
type Attachment struct {
	ID   uuid.UUID `json:"id"`
	URL  string    `json:"url"`
	Mime string    `json:"mime"`
	Name string    `json:"name"`
	Size int64     `json:"size"`
}

// ChatSummary is a conversation list item. TaskID is the public task id so the
// chat UI and the task pages share one identifier scheme.
type ChatSummary struct {
	ID        string    `json:"id"`
	TaskID    string    `json:"taskId"`
	Name      string    `json:"name"`
	LastMessage string  `json:"last"`
	Unread    int       `json:"unread"`
	Online    bool      `json:"online"`
	UpdatedAt time.Time `json:"time"`
}

// SendMessageInput holds message creation data. The receiver is derived from the
// task (the other participant), so callers do not pass it.
type SendMessageInput struct {
	// TaskRef is the task's public id (e.g. "TB-7") or UUID; resolved by the service.
	TaskRef    string
	SenderID   uuid.UUID
	Message    string
	Attachment *Attachment
}

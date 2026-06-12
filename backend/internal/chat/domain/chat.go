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

// ChatSummary is a conversation list item.
type ChatSummary struct {
	ID        uuid.UUID `json:"id"`
	TaskID    uuid.UUID `json:"taskId"`
	Name      string    `json:"name"`
	LastMessage string  `json:"last"`
	Unread    int       `json:"unread"`
	Online    bool      `json:"online"`
	UpdatedAt time.Time `json:"time"`
}

// SendMessageInput holds message creation data.
type SendMessageInput struct {
	TaskID      uuid.UUID
	SenderID    uuid.UUID
	ReceiverID  uuid.UUID
	Message     string
	Attachment  *Attachment
}

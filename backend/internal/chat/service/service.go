package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/chat/domain"
)

// Service defines chat use cases. taskRef is the task's public id or UUID.
type Service interface {
	ListChats(ctx context.Context, userID uuid.UUID) ([]domain.ChatSummary, error)
	ListMessages(ctx context.Context, taskRef string, userID uuid.UUID, before *uuid.UUID, limit int) ([]domain.ChatMessage, error)
	SendMessage(ctx context.Context, input domain.SendMessageInput) (*domain.ChatMessage, error)
	MarkRead(ctx context.Context, taskRef string, userID uuid.UUID) error
}

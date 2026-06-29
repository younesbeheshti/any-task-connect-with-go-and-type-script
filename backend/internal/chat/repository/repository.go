package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/chat/domain"
)

// Repository defines chat persistence operations.
type Repository interface {
	CreateMessage(ctx context.Context, msg *domain.ChatMessage) error
	GetMessageByID(ctx context.Context, id uuid.UUID) (*domain.ChatMessage, error)
	ListMessages(ctx context.Context, taskID uuid.UUID, before *uuid.UUID, limit int) ([]domain.ChatMessage, error)
	MarkSeen(ctx context.Context, taskID, receiverID uuid.UUID) error
	CountUnread(ctx context.Context, userID uuid.UUID) (int, error)
	ListChatsForUser(ctx context.Context, userID uuid.UUID) ([]domain.ChatSummary, error)
	// LatestCounterparty returns the other participant of the user's most recent
	// message in a task, if any (used to resolve the reply recipient).
	LatestCounterparty(ctx context.Context, taskID, userID uuid.UUID) (uuid.UUID, bool, error)
}

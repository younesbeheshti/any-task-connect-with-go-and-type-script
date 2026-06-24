package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/chat/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/chat/repository"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
)

type ChatService struct {
	repo repository.Repository
}

func NewChatService(repo repository.Repository) *ChatService {
	return &ChatService{repo: repo}
}

func (s *ChatService) ListChats(ctx context.Context, userID uuid.UUID) ([]domain.ChatSummary, error) {
	return s.repo.ListChatsForUser(ctx, userID)
}

func (s *ChatService) ListMessages(ctx context.Context, taskID, userID uuid.UUID, before *uuid.UUID, limit int) ([]domain.ChatMessage, error) {
	msgs, err := s.repo.ListMessages(ctx, taskID, before, limit)
	if err != nil {
		return nil, err
	}
	// mark received messages as seen
	_ = s.repo.MarkSeen(ctx, taskID, userID)
	return msgs, nil
}

func (s *ChatService) SendMessage(ctx context.Context, input domain.SendMessageInput) (*domain.ChatMessage, error) {
	if input.Message == "" && input.Attachment == nil {
		return nil, apperrors.New("INVALID_INPUT", "پیام یا پیوست الزامی است", 422, apperrors.ErrValidation)
	}
	msg := &domain.ChatMessage{
		ID:         uuid.New(),
		TaskID:     input.TaskID,
		SenderID:   input.SenderID,
		ReceiverID: input.ReceiverID,
		Message:    input.Message,
		Attachment: input.Attachment,
		Seen:       false,
	}
	if err := s.repo.CreateMessage(ctx, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (s *ChatService) MarkRead(ctx context.Context, taskID, userID uuid.UUID) error {
	return s.repo.MarkSeen(ctx, taskID, userID)
}

func (s *ChatService) CountUnread(ctx context.Context, userID uuid.UUID) (int, error) {
	return s.repo.CountUnread(ctx, userID)
}

// Ensure *ChatService satisfies Service + common unread counter.
var _ Service = (*ChatService)(nil)

// UnreadCounter extends the service with unread count.
type UnreadCounter interface {
	CountUnread(ctx context.Context, userID uuid.UUID) (int, error)
}

var _ UnreadCounter = (*ChatService)(nil)


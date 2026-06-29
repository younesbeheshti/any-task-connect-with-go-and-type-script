package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/chat/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/chat/repository"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	taskrepo "github.com/younesbeheshti/any-task-connect/backend/internal/task/repository"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
)

type ChatService struct {
	repo  repository.Repository
	tasks taskrepo.Repository
}

func NewChatService(repo repository.Repository, tasks taskrepo.Repository) *ChatService {
	return &ChatService{repo: repo, tasks: tasks}
}

// resolveTaskID accepts a task public id (e.g. "TB-7") or a UUID and returns the
// internal task UUID along with the task's requester and assigned agent.
func (s *ChatService) resolveTask(ctx context.Context, ref string) (taskID, requesterID uuid.UUID, agentID *uuid.UUID, err error) {
	if id, e := uuid.Parse(ref); e == nil {
		t, e2 := s.tasks.GetByID(ctx, id)
		if e2 == nil {
			return t.ID, t.RequesterID, t.AssignedAgentID, nil
		}
		if !errors.Is(e2, common.ErrNotFound) {
			return uuid.Nil, uuid.Nil, nil, e2
		}
	}
	t, e := s.tasks.GetByPublicID(ctx, ref)
	if errors.Is(e, common.ErrNotFound) {
		return uuid.Nil, uuid.Nil, nil, apperrors.New("NOT_FOUND", "تسک یافت نشد", 404, apperrors.ErrNotFound)
	}
	if e != nil {
		return uuid.Nil, uuid.Nil, nil, e
	}
	return t.ID, t.RequesterID, t.AssignedAgentID, nil
}

func (s *ChatService) ListChats(ctx context.Context, userID uuid.UUID) ([]domain.ChatSummary, error) {
	return s.repo.ListChatsForUser(ctx, userID)
}

func (s *ChatService) ListMessages(ctx context.Context, taskRef string, userID uuid.UUID, before *uuid.UUID, limit int) ([]domain.ChatMessage, error) {
	taskID, _, _, err := s.resolveTask(ctx, taskRef)
	if err != nil {
		return nil, err
	}
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
	taskID, requesterID, agentID, err := s.resolveTask(ctx, input.TaskRef)
	if err != nil {
		return nil, err
	}

	receiverID, err := s.resolveReceiver(ctx, taskID, input.SenderID, requesterID, agentID)
	if err != nil {
		return nil, err
	}

	msg := &domain.ChatMessage{
		ID:         uuid.New(),
		TaskID:     taskID,
		SenderID:   input.SenderID,
		ReceiverID: receiverID,
		Message:    input.Message,
		Attachment: input.Attachment,
		Seen:       false,
	}
	if err := s.repo.CreateMessage(ctx, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

// resolveReceiver picks the other participant: continue an existing thread with the
// last counterparty, otherwise fall back to the task roles (agent ↔ requester).
func (s *ChatService) resolveReceiver(ctx context.Context, taskID, senderID, requesterID uuid.UUID, agentID *uuid.UUID) (uuid.UUID, error) {
	if peer, ok, err := s.repo.LatestCounterparty(ctx, taskID, senderID); err == nil && ok {
		return peer, nil
	}
	if senderID == requesterID {
		if agentID != nil {
			return *agentID, nil
		}
		return uuid.Nil, apperrors.New("NO_RECIPIENT", "هنوز مجری‌ای برای گفت‌وگو انتخاب نشده است", 409, apperrors.ErrConflict)
	}
	return requesterID, nil
}

func (s *ChatService) MarkRead(ctx context.Context, taskRef string, userID uuid.UUID) error {
	taskID, _, _, err := s.resolveTask(ctx, taskRef)
	if err != nil {
		return err
	}
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


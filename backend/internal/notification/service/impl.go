package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/notification/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/notification/repository"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
)

type NotificationService struct {
	repo repository.Repository
}

func NewNotificationService(repo repository.Repository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) Create(ctx context.Context, input domain.CreateNotificationInput) (*domain.Notification, error) {
	n := &domain.Notification{
		ID:       uuid.New(),
		UserID:   input.UserID,
		Title:    input.Title,
		Body:     input.Body,
		Type:     input.Type,
		IsRead:   false,
		DeepLink: input.DeepLink,
	}
	if err := s.repo.Create(ctx, n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *NotificationService) GetByID(ctx context.Context, id, userID uuid.UUID) (*domain.Notification, error) {
	n, err := s.repo.GetByID(ctx, id)
	if err == common.ErrNotFound {
		return nil, apperrors.New("NOT_FOUND", "اعلان یافت نشد", 404, apperrors.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	if n.UserID != userID {
		return nil, apperrors.New("FORBIDDEN", "دسترسی مجاز نیست", 403, apperrors.ErrForbidden)
	}
	return n, nil
}

func (s *NotificationService) List(ctx context.Context, userID uuid.UUID, unreadOnly bool, pg common.PaginationParams) (*common.PaginatedResult[domain.Notification], error) {
	items, total, err := s.repo.ListByUser(ctx, userID, unreadOnly, pg)
	if err != nil {
		return nil, err
	}
	page := pg.Page
	if page < 1 {
		page = 1
	}
	pageSize := pg.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	return &common.PaginatedResult[domain.Notification]{
		Items:    items,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

func (s *NotificationService) MarkRead(ctx context.Context, id, userID uuid.UUID) error {
	err := s.repo.MarkRead(ctx, id, userID)
	if err == common.ErrNotFound {
		return apperrors.New("NOT_FOUND", "اعلان یافت نشد", 404, apperrors.ErrNotFound)
	}
	return err
}

func (s *NotificationService) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	return s.repo.MarkAllRead(ctx, userID)
}

func (s *NotificationService) CountUnread(ctx context.Context, userID uuid.UUID) (int64, error) {
	return s.repo.CountUnread(ctx, userID)
}

var _ Service = (*NotificationService)(nil)

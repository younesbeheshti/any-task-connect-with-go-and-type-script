package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/rating/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/rating/repository"
	apperrors "github.com/younesbeheshti/any-task-connect/backend/pkg/errors"
)

type ReviewService struct {
	repo repository.Repository
}

func NewReviewService(repo repository.Repository) *ReviewService {
	return &ReviewService{repo: repo}
}

func (s *ReviewService) Create(ctx context.Context, input domain.CreateReviewInput) (*domain.Review, error) {
	if input.Rating < 1 || input.Rating > 5 {
		return nil, apperrors.New("INVALID_RATING", "امتیاز باید بین ۱ تا ۵ باشد", 422, apperrors.ErrValidation)
	}
	// prevent duplicate review
	if _, err := s.repo.GetByTaskAndReviewer(ctx, input.TaskID, input.ReviewerID); err == nil {
		return nil, apperrors.New("DUPLICATE_REVIEW", "شما قبلاً برای این درخواست نظر ثبت کرده‌اید", 409, apperrors.ErrConflict)
	}
	review := &domain.Review{
		ID:             uuid.New(),
		TaskID:         input.TaskID,
		ReviewerID:     input.ReviewerID,
		ReviewedUserID: input.ReviewedUserID,
		Rating:         input.Rating,
		Comment:        input.Comment,
	}
	if err := s.repo.Create(ctx, review); err != nil {
		return nil, err
	}
	return review, nil
}

func (s *ReviewService) GetByTask(ctx context.Context, taskID uuid.UUID) ([]domain.Review, error) {
	return s.repo.ListByTask(ctx, taskID)
}

func (s *ReviewService) GetByUser(ctx context.Context, userID uuid.UUID) ([]domain.Review, error) {
	return s.repo.ListByReviewedUser(ctx, userID)
}

// AverageRating returns avg + count for a user.
func (s *ReviewService) AverageRating(ctx context.Context, userID uuid.UUID) (float64, int, error) {
	return s.repo.AverageRating(ctx, userID)
}

var _ Service = (*ReviewService)(nil)


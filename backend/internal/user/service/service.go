package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/user/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/user/repository"
)

// Service defines user profile use cases.
type Service interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetPublicProfile(ctx context.Context, id uuid.UUID) (*domain.PublicProfile, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, input domain.UpdateProfileInput) (*domain.User, error)
	ListUsers(ctx context.Context, filter repository.UserFilter, pg common.PaginationParams) (*common.PaginatedResult[domain.User], error)
	Deactivate(ctx context.Context, userID uuid.UUID) error
	UpdateRatingAggregate(ctx context.Context, userID uuid.UUID) error
}

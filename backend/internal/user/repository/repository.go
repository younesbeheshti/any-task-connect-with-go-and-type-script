package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/internal/user/domain"
)

// Repository defines user persistence operations.
type Repository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByPhone(ctx context.Context, phone string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error
	UpdateRating(ctx context.Context, userID uuid.UUID, rating float64, count int) error
	IncrementCompletedTasks(ctx context.Context, userID uuid.UUID) error
	List(ctx context.Context, filter UserFilter, pg common.PaginationParams) ([]domain.User, int64, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

// UserFilter supports admin user search.
type UserFilter struct {
	Query    string
	Role     *common.Role
	IsActive *bool
}

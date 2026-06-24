package repository

import (
	"context"

	"github.com/younesbeheshti/any-task-connect/backend/internal/audit/domain"
)

type Repository interface {
	Create(ctx context.Context, log *domain.AuditLog) error
}

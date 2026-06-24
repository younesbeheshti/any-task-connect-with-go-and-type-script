package service

import (
	"context"

	"github.com/younesbeheshti/any-task-connect/backend/internal/audit/domain"
)

type Service interface {
	Log(ctx context.Context, input domain.CreateInput) error
}

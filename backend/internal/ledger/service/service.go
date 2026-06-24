package service

import (
	"context"

	"github.com/younesbeheshti/any-task-connect/backend/internal/ledger/domain"
)

type Service interface {
	Record(ctx context.Context, entries []domain.Entry) error
}

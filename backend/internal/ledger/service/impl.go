package service

import (
	"context"
	"fmt"

	"github.com/younesbeheshti/any-task-connect/backend/internal/ledger/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/ledger/repository"
)

type LedgerService struct {
	repo repository.Repository
}

func NewLedgerService(repo repository.Repository) *LedgerService {
	return &LedgerService{repo: repo}
}

func (s *LedgerService) Record(ctx context.Context, entries []domain.Entry) error {
	if !domain.ValidateDoubleEntry(entries) {
		return fmt.Errorf("double-entry validation failed: debits must equal credits")
	}
	return s.repo.Create(ctx, entries)
}

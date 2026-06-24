package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/audit/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/audit/repository"
)

type AuditService struct {
	repo repository.Repository
}

func NewAuditService(repo repository.Repository) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) Log(ctx context.Context, input domain.CreateInput) error {
	log := &domain.AuditLog{
		ID:          uuid.New(),
		EntityType:  input.EntityType,
		EntityID:    input.EntityID,
		Action:      input.Action,
		PerformedBy: input.PerformedBy,
		OldValue:    input.OldValue,
		NewValue:    input.NewValue,
	}
	return s.repo.Create(ctx, log)
}

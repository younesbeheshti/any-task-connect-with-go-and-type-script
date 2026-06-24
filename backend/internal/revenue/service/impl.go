package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/revenue/domain"
	"github.com/younesbeheshti/any-task-connect/backend/internal/revenue/repository"
)

type RevenueService struct {
	repo repository.Repository
}

func NewRevenueService(repo repository.Repository) *RevenueService {
	return &RevenueService{repo: repo}
}

func (s *RevenueService) RecordCommission(ctx context.Context, taskID *uuid.UUID, amount int64) error {
	rev := &domain.Revenue{
		TaskID: taskID,
		Source: "commission",
		Amount: amount,
	}
	return s.repo.Create(ctx, rev)
}

func (s *RevenueService) GetStats(ctx context.Context) (*domain.RevenueStats, error) {
	return s.repo.GetStats(ctx)
}

func (s *RevenueService) GetDaily(ctx context.Context, days int) ([]domain.DailyRevenue, error) {
	if days <= 0 {
		days = 30
	}
	return s.repo.GetDaily(ctx, days)
}

func (s *RevenueService) GetMonthly(ctx context.Context, months int) ([]domain.MonthlyRevenue, error) {
	if months <= 0 {
		months = 12
	}
	return s.repo.GetMonthly(ctx, months)
}

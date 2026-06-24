package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/younesbeheshti/any-task-connect/backend/internal/revenue/domain"
)

type Service interface {
	RecordCommission(ctx context.Context, taskID *uuid.UUID, amount int64) error
	GetStats(ctx context.Context) (*domain.RevenueStats, error)
	GetDaily(ctx context.Context, days int) ([]domain.DailyRevenue, error)
	GetMonthly(ctx context.Context, months int) ([]domain.MonthlyRevenue, error)
}

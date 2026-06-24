package repository

import (
	"context"

	"github.com/younesbeheshti/any-task-connect/backend/internal/revenue/domain"
)

type Repository interface {
	Create(ctx context.Context, r *domain.Revenue) error
	GetStats(ctx context.Context) (*domain.RevenueStats, error)
	GetDaily(ctx context.Context, days int) ([]domain.DailyRevenue, error)
	GetMonthly(ctx context.Context, months int) ([]domain.MonthlyRevenue, error)
	GetTotal(ctx context.Context) (int64, error)
	GetMonthToDate(ctx context.Context) (int64, error)
}

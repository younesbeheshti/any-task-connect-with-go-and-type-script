package domain

import (
	"time"

	"github.com/google/uuid"
)

type Revenue struct {
	ID        uuid.UUID
	TaskID    *uuid.UUID
	Source    string
	Amount    int64
	CreatedAt time.Time
}

type DailyRevenue struct {
	Date   string `json:"date"`
	Amount int64  `json:"amount"`
}

type MonthlyRevenue struct {
	Month  string `json:"month"`
	Amount int64  `json:"amount"`
}

type RevenueStats struct {
	TotalRevenue   int64            `json:"totalRevenue"`
	MonthToDate    int64            `json:"monthToDate"`
	DailyRevenue   []DailyRevenue   `json:"dailyRevenue"`
	MonthlyRevenue []MonthlyRevenue `json:"monthlyRevenue"`
}

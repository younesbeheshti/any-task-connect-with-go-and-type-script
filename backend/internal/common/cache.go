package common

import (
	"context"
	"time"
)

// Cache abstracts Redis operations used across the application.
type Cache interface {
	Get(ctx context.Context, key string, dest any) error
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, key string) (bool, error)
	Incr(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, ttl time.Duration) error
}

// Cache key builders.
const (
	KeyJWTBlacklist   = "jwt:blacklist:"
	KeySession        = "session:"
	KeyRateLimit      = "ratelimit:"
	KeyOnlineUser     = "online:"
	KeyTaskCache      = "task:"
	KeyNotification   = "notifications:"
	KeyWalletCache    = "wallet:"
	KeyWalletStats    = "wallet_stats:"
	KeyRevenueDaily   = "revenue:daily"
	KeyRevenueMonthly = "revenue:monthly"
)

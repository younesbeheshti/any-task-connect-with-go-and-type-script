package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/younesbeheshti/any-task-connect/backend/configs"
	"go.uber.org/zap"
)

// Client wraps the Redis client.
type Client struct {
	*redis.Client
	log *zap.Logger
}

// Connect establishes a Redis client and verifies connectivity.
func Connect(ctx context.Context, cfg configs.RedisConfig, log *zap.Logger) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	log.Info("redis connected", zap.String("addr", cfg.Addr()))

	return &Client{Client: rdb, log: log}, nil
}

// Ping checks Redis connectivity.
func (c *Client) Ping(ctx context.Context) error {
	return c.Client.Ping(ctx).Err()
}

// HealthCheck verifies Redis is reachable.
func (c *Client) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return c.Ping(ctx)
}

// Close closes the Redis connection.
func (c *Client) Close() error {
	return c.Client.Close()
}

package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimit applies a sliding-window rate limit per IP and route.
func RateLimit(rdb *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	if limit <= 0 {
		return func(c *gin.Context) { c.Next() }
	}

	return func(c *gin.Context) {
		key := fmt.Sprintf("ratelimit:%s:%s", c.ClientIP(), c.FullPath())
		ctx := c.Request.Context()

		allowed, err := slidingWindowAllow(ctx, rdb, key, limit, window)
		if err != nil {
			c.Next()
			return
		}
		if !allowed {
			c.AbortWithStatusJSON(429, gin.H{
				"success": false,
				"message": "rate limit exceeded",
				"errors":  gin.H{"code": "RATE_LIMIT_EXCEEDED"},
			})
			return
		}
		c.Next()
	}
}

func slidingWindowAllow(ctx context.Context, rdb *redis.Client, key string, limit int, window time.Duration) (bool, error) {
	now := time.Now().UTC()
	windowStart := now.Add(-window)

	pipe := rdb.Pipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart.UnixNano()))
	countCmd := pipe.ZCard(ctx, key)
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(now.UnixNano()), Member: now.UnixNano()})
	pipe.Expire(ctx, key, window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return true, err
	}

	return countCmd.Val() < int64(limit), nil
}

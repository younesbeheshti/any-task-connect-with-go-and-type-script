package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
)

// Cache implements common.Cache using Redis.
type Cache struct {
	client *goredis.Client
}

// NewCache creates a Redis-backed cache.
func NewCache(client *goredis.Client) *Cache {
	return &Cache{client: client}
}

func (c *Cache) Get(ctx context.Context, key string, dest any) error {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (c *Cache) GetString(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *Cache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *Cache) SetString(ctx context.Context, key, value string, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *Cache) Delete(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	return n > 0, err
}

func (c *Cache) Incr(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, key).Result()
}

func (c *Cache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return c.client.Expire(ctx, key, ttl).Err()
}

// BlacklistJTI adds a JWT ID to the revocation blacklist.
func (c *Cache) BlacklistJTI(ctx context.Context, jti string, ttl time.Duration) error {
	key := common.KeyJWTBlacklist + jti
	return c.SetString(ctx, key, "1", ttl)
}

// IsJTIBlacklisted reports whether a JWT ID is revoked.
func (c *Cache) IsJTIBlacklisted(ctx context.Context, jti string) (bool, error) {
	key := common.KeyJWTBlacklist + jti
	return c.Exists(ctx, key)
}

// SessionKey builds a session Redis key.
func SessionKey(userID, sessionID string) string {
	return fmt.Sprintf("%s%s:%s", common.KeySession, userID, sessionID)
}

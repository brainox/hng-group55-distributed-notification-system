package idempotency

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Checker struct {
	redis *redis.Client
	ttl   time.Duration
}

func NewChecker(redis *redis.Client, ttl time.Duration) *Checker {
	return &Checker{
		redis: redis,
		ttl:   ttl,
	}
}

// IsProcessed checks if a message has already been processed
func (c *Checker) IsProcessed(ctx context.Context, messageID string) (bool, error) {
	key := fmt.Sprintf("email:processed:%s", messageID)
	exists, err := c.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// MarkProcessed marks a message as processed
func (c *Checker) MarkProcessed(ctx context.Context, messageID string) error {
	key := fmt.Sprintf("email:processed:%s", messageID)
	return c.redis.Set(ctx, key, "1", c.ttl).Err()
}

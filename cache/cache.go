package cache

import (
	"context"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string) ([]byte, bool, error)
	Del(ctx context.Context, key string) error
	Expire(ctx context.Context, key string, ttl time.Duration) error

	SetNX(ctx context.Context, key, value string, ttl time.Duration) (bool, error)
	Incr(ctx context.Context, key string, ttl time.Duration) (int64, error)

	Ping(ctx context.Context) error
	Close() error
}

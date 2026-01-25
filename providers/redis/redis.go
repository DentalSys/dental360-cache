package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr         string
	Password     string
	DB           int
	Prefix       string
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Cache struct {
	rdb    *redis.Client
	prefix string
}

func New(cfg Config) (*Cache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	c := &Cache{rdb: rdb, prefix: cfg.Prefix}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := c.Ping(ctx); err != nil {
		_ = rdb.Close()
		return nil, err
	}

	return c, nil
}

func (c *Cache) key(k string) string {
	if c.prefix == "" {
		return k
	}
	return c.prefix + k
}

func (c *Cache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	switch v := value.(type) {
	case []byte:
		return c.rdb.Set(ctx, c.key(key), v, ttl).Err()
	case string:
		return c.rdb.Set(ctx, c.key(key), v, ttl).Err()
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		return c.rdb.Set(ctx, c.key(key), b, ttl).Err()
	}
}

func (c *Cache) Get(ctx context.Context, key string) ([]byte, bool, error) {
	v, err := c.rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return v, true, nil
}

func (c *Cache) Del(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, c.key(key)).Err()
}

func (c *Cache) SetNX(ctx context.Context, key, value string, ttl time.Duration) (bool, error) {
	return c.rdb.SetNX(ctx, key, value, ttl).Result()
}

func (c *Cache) Incr(ctx context.Context, key string, ttl time.Duration) (int64, error) {
	n, err := c.rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, nil
	}

	if n == 1 && ttl > 0 {
		_ = c.rdb.Expire(ctx, key, ttl).Err()
	}

	return n, nil
}

func (c *Cache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return c.rdb.Expire(ctx, c.key(key), ttl).Err()
}

func (c *Cache) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

func (c *Cache) Close() error {
	return c.rdb.Close()
}

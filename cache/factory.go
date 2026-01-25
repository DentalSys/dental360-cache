package cache

import (
	"fmt"
	"time"

	"github.com/Dentalsys/dental360-cache/providers/redis"
)

func New(cfg Config) (Cache, error) {
	switch cfg.Provider {
	case ProviderRedis:
		if cfg.Redis.Addr == "" {
			return nil, fmt.Errorf("%w: redis.addr is required", ErrInvalidConfig)
		}

		normalizeRedis(&cfg)
		return redis.New(redis.Config{
			Addr:         cfg.Redis.Addr,
			Password:     cfg.Redis.Password,
			DB:           cfg.Redis.DB,
			Prefix:       cfg.Prefix,
			DialTimeout:  cfg.Redis.DialTimeout,
			ReadTimeout:  cfg.Redis.ReadTimeout,
			WriteTimeout: cfg.Redis.WriteTimeout,
		})
	default:
		return nil, fmt.Errorf("%w: unsupported provider %q", ErrInvalidConfig, cfg.Provider)
	}
}

func normalizeRedis(cfg *Config) {
	if cfg.Redis.DialTimeout <= 0 {
		cfg.Redis.DialTimeout = 3 * time.Second
	}
	if cfg.Redis.ReadTimeout <= 0 {
		cfg.Redis.ReadTimeout = 2 * time.Second
	}
	if cfg.Redis.WriteTimeout <= 0 {
		cfg.Redis.WriteTimeout = 2 * time.Second
	}
}

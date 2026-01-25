package cache

import "time"

type Provider string

const (
	ProviderRedis Provider = "redis"
)

type Config struct {
	Provider Provider
	Prefix   string
	Redis    RedisConfig
}

type RedisConfig struct {
	Addr         string
	Password     string
	DB           int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

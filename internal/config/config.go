package cfg

import (
	"time"

	"github.com/j3dyy/nazuki/internal/env"
)

type Option func(*Config)

type dbConfig struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

type redisConfig struct {
	addr     string
	password string
}

type natsConfig struct {
	url string
}

type Config struct {
	addr        string
	DBConfig    dbConfig
	RedisConfig redisConfig
	NatsConfig  natsConfig
}

func WithAddr(addr string) Option {
	return func(c *Config) {
		c.addr = addr
	}
}

func WithDBConfig(dsn string, maxOpenConns, maxIdleConns int, maxIdleTime time.Duration) Option {
	return func(c *Config) {
		c.DBConfig = dbConfig{
			DSN:          dsn,
			MaxOpenConns: maxOpenConns,
			MaxIdleConns: maxIdleConns,
			MaxIdleTime:  maxIdleTime,
		}
	}
}

func WithRedisConfig(addr, password string) Option {
	return func(c *Config) {
		c.RedisConfig = redisConfig{
			addr:     addr,
			password: password,
		}
	}
}

func WithNatsConfig(url string) Option {
	return func(c *Config) {
		c.NatsConfig = natsConfig{
			url: url,
		}
	}
}

func NewConfig(opts ...Option) *Config {
	c := &Config{
		addr:        "localhost:8080",
		DBConfig:    dbConfig{},
		RedisConfig: redisConfig{},
		NatsConfig:  natsConfig{},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func LoadConfigFromEnv() *Config {
	defaultCfg := NewConfig(
		WithAddr(env.GetString("APP_ADDR", ":8080")),
		WithDBConfig(
			env.GetString("DB_DSN", ""),
			env.GetInt("MAX_DB_OPEN_CONNS", 30),
			env.GetInt("MAX_DB_IDLE_CONNS", 30),
			env.GetDuration("MAX_DB_IDLE_TIME", time.Second*30),
		),
		WithRedisConfig(
			env.GetString("REDIS_ADDR", "localhost"),
			env.GetString("REDIS_PASSWORD", ""),
		),
		WithNatsConfig(env.GetString("NATS_URL", "nats://localhost:4222")),
	)
	return defaultCfg
}

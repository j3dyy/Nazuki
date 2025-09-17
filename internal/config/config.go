package cfg

import (
	"os"
	"time"

	"github.com/j3dyy/nazuki/internal/env"
	"github.com/j3dyy/nazuki/internal/service"
	"github.com/j3dyy/nazuki/internal/store"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type Option func(*config)

type dbConfig struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  time.Duration
}

type redisConfig struct {
	addr     string
	password string
}

type natsConfig struct {
	url string
}

type config struct {
	addr        string
	dbConfig    dbConfig
	redisConfig redisConfig
	natsConfig  natsConfig
}

func WithAddr(addr string) Option {
	return func(c *config) {
		c.addr = addr
	}
}

func WithDBConfig(dsn string, maxOpenConns, maxIdleConns int, maxIdleTime time.Duration) Option {
	return func(c *config) {
		c.dbConfig = dbConfig{
			dsn:          dsn,
			maxOpenConns: maxOpenConns,
			maxIdleConns: maxIdleConns,
			maxIdleTime:  maxIdleTime,
		}
	}
}

func WithRedisConfig(addr, password string) Option {
	return func(c *config) {
		c.redisConfig = redisConfig{
			addr:     addr,
			password: password,
		}
	}
}

func WithNatsConfig(url string) Option {
	return func(c *config) {
		c.natsConfig = natsConfig{
			url: url,
		}
	}
}

func NewConfig(opts ...Option) *config {
	c := &config{
		addr:        "localhost:8080",
		dbConfig:    dbConfig{},
		redisConfig: redisConfig{},
		natsConfig:  natsConfig{},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func LoadConfigFromEnv() *config {
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

type Version string

type Application struct {
	Version     Version
	Logger      zerolog.Logger
	Service     service.Service
	Store       store.Store
	RedisClient *redis.Client
	NatsClient  *nats.Conn
}

func NewApplication(service service.Service, store store.Store) Application {
	app := Application{
		Version: Version("0.0.1"),
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}
	app.Logger = zerolog.New(output).With().Timestamp().Logger()

	return app
}

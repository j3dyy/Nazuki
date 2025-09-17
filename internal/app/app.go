package app

import (
	"os"
	"sync"
	"time"

	cfg "github.com/j3dyy/nazuki/internal/config"
	"github.com/j3dyy/nazuki/internal/db"
	"github.com/j3dyy/nazuki/internal/service"
	"github.com/j3dyy/nazuki/internal/store"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type Version string

type Application struct {
	Version       Version
	Logger        zerolog.Logger
	Wait          *sync.WaitGroup
	ErrorChan     chan error
	ErrorChanDone chan bool
	Service       service.Service
	Store         store.Store
	RedisClient   *redis.Client
	NatsClient    *nats.Conn
}

func NewApplication(cfg *cfg.Config) (*Application, error) {

	db, err := db.NewPostgres(cfg.DBConfig.DSN, cfg.DBConfig.MaxOpenConns, cfg.DBConfig.MaxIdleConns, cfg.DBConfig.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	store := store.NewStore(db)
	service := service.NewService(store)

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	app := &Application{
		Version:       Version("0.0.1"),
		Logger:        zerolog.New(output).With().Timestamp().Logger(),
		Wait:          &sync.WaitGroup{},
		ErrorChan:     make(chan error),
		ErrorChanDone: make(chan bool),
		Service:       service,
		Store:         *store,
	}

	return app, nil
}

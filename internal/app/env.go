package app

import (
	"context"
	"log/slog"

	"github.com/Sanchir01/avito-testovoe/internal/config"
)

type Env struct {
	Lg       *slog.Logger
	Cfg      *config.Config
	Database *Database
	Handlers *Handlers
}

func NewEnv() (*Env, error) {
	cfg := config.MustLoadConfig()
	logger := setupLogger(cfg.Env)
	ctx := context.Background()
	primarydb, err := NewDataBases(ctx, cfg.PrimaryDB.User, cfg.PrimaryDB.Host, cfg.PrimaryDB.Dbname, cfg.PrimaryDB.Port, cfg.PrimaryDB.MaxAttempts)
	if err != nil {
		return nil, err
	}
	repositories := NewRepositories(primarydb)
	services := NewServices(repositories, primarydb)
	handlers := NewHandlers(services, logger)
	env := Env{
		Lg:       logger,
		Cfg:      cfg,
		Database: primarydb,
		Handlers: handlers,
	}

	return &env, nil
}

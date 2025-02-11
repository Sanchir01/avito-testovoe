package app

import (
	"context"
	"github.com/Sanchir01/avito-testovoe/internal/config"
	"log/slog"
)

type Env struct {
	Lg       *slog.Logger
	Cfg      *config.Config
	Services *Services
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
	services := NewServices(repositories)
	env := Env{
		Lg:       logger,
		Cfg:      cfg,
		Services: services,
	}

	return &env, nil
}

package app

import (
	"github.com/Sanchir01/avito-testovoe/internal/config"
	"log/slog"
)

type Env struct {
	Lg  *slog.Logger
	Cfg *config.Config
}

func NewEnv() (*Env, error) {
	cfg := config.MustLoadConfig()
	logger := setupLogger(cfg.Env)
	primarydb, err := NewDataBases(ctx,cfg.)
	env := Env{
		Lg:  logger,
		Cfg: cfg,
	}
	return &env, nil
}

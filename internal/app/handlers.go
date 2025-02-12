package app

import (
	"github.com/Sanchir01/avito-testovoe/internal/feature/user"
	"log/slog"
)

type Handlers struct {
	UserHandler *user.Handler
}

func NewHandlers(services *Services, log *slog.Logger) *Handlers {
	return &Handlers{
		UserHandler: user.NewHandler(services.UserService, log),
	}
}

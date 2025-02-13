package app

import (
	"log/slog"

	"github.com/Sanchir01/avito-testovoe/internal/feature/user"
)

type Handlers struct {
	UserHandler *user.Handler
}

func NewHandlers(services *Services, log *slog.Logger) *Handlers {
	return &Handlers{
		UserHandler: user.NewHandler(services.UserService, log),
	}
}

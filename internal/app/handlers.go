package app

import (
	"github.com/Sanchir01/avito-testovoe/internal/feature/product"
	"log/slog"

	"github.com/Sanchir01/avito-testovoe/internal/feature/user"
)

type Handlers struct {
	UserHandler    *user.Handler
	ProductHandler *product.Handler
}

func NewHandlers(services *Services, log *slog.Logger) *Handlers {
	return &Handlers{
		UserHandler:    user.NewHandler(services.UserService, log),
		ProductHandler: product.NewHandler(services.ProductService, log),
	}
}

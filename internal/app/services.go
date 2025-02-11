package app

import (
	"github.com/Sanchir01/avito-testovoe/internal/feature/product"
	"github.com/Sanchir01/avito-testovoe/internal/feature/user"
)

type Services struct {
	ProductService *product.Service
	UserService    *user.Service
}

func NewServices(r *Repositories, db *Database) *Services {
	return &Services{
		ProductService: product.NewService(r.ProductRepository),
		UserService:    user.NewService(r.UserRepository, db.PrimaryDB),
	}
}

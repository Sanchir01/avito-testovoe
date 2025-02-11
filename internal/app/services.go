package app

import "github.com/Sanchir01/avito-testovoe/internal/feature/product"

type Services struct {
	ProductService *product.Service
}

func NewServices(r *Repositories) *Services {
	return &Services{
		ProductService: product.NewService(r.ProductRepository),
	}
}

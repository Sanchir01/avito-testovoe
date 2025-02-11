package app

import (
	"github.com/Sanchir01/avito-testovoe/internal/feature/product"
	"github.com/Sanchir01/avito-testovoe/internal/feature/user"
)

type Repositories struct {
	ProductRepository *product.Repository
	UserRepository    *user.Repository
}

func NewRepositories(database *Database) *Repositories {
	return &Repositories{
		ProductRepository: product.NewRepository(database.PrimaryDB),
		UserRepository:    user.NewRepository(database.PrimaryDB),
	}
}

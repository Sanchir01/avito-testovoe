package app

import "github.com/Sanchir01/avito-testovoe/internal/feature/product"

type Repositories struct {
	ProductRepository *product.Repository
}

func NewRepositories(database *Database) *Repositories {
	return &Repositories{
		ProductRepository: product.NewRepository(database.PrimaryDB),
	}
}

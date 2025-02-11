package product

import "github.com/jackc/pgx/v5/pgxpool"

type Repository struct {
	primaryDB *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{primaryDB: db}
}

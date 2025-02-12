package product

import (
	"github.com/google/uuid"
)

type DataBaseProduct struct {
	ID      uuid.UUID `db:"id"`
	Title   string    `db:"title"`
	Slug    string    `db:"slug"`
	Version int64     `db:"version"`
	Price   int64     `db:"price"`
}

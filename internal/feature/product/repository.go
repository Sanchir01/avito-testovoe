package product

import (
	"context"
	"fmt"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	primaryDB *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{primaryDB: db}
}

func (r *Repository) GetProductByID(ctx context.Context, id uuid.UUID) (*DataBaseProduct, error) {
	conn, err := r.primaryDB.Acquire(ctx)
	if err != nil {
		slog.Error("Failed to acquire DB connection", slog.String("error", err.Error()))
		return nil, err
	}
	defer conn.Release()

	query, args, err := sq.
		Select("id, title,slug,version,price").
		From("product").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	var products DataBaseProduct

	if err := conn.QueryRow(ctx, query, args...).Scan(&products.ID, &products.Title, &products.Slug,
		&products.Version, &products.Price); err != nil {
		return nil, err
	}
	fmt.Println("product repo", products)
	return &products, nil
}
func (r *Repository) GetAllProducts(ctx context.Context) ([]*DataBaseProduct, error) {
	conn, err := r.primaryDB.Acquire(ctx)
	if err != nil {
		slog.Error("Failed to acquire DB connection", slog.String("error", err.Error()))
		return nil, err
	}
	defer conn.Release()

	query, args, err := sq.
		Select("id, title,slug,version,price").
		From("product").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []*DataBaseProduct
	for rows.Next() {
		var product DataBaseProduct
		if err := rows.Scan(&product.ID, &product.Title, &product.Slug, &product.Version, &product.Price); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	return products, nil
}

package product

import (
	"context"
	"fmt"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	primaryDB *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{primaryDB: db}
}
func (r *Repository) CreateProduct(ctx context.Context, title, slug string, price int, tx pgx.Tx) error {
	query, arg, err := sq.Insert("product").
		Columns("title", "slug", "price").
		Values(title, slug, price).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}
	_, exceErr := tx.Exec(ctx, query, arg...)
	if exceErr != nil {
		return fmt.Errorf("failed to create product: %s", exceErr.Error())
	}
	return nil
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
	slog.Error("query", query)
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

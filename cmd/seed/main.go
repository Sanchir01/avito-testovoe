package main

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/Sanchir01/avito-testovoe/internal/app"
	"github.com/jackc/pgx/v5"
)

type Product struct {
	Title string
	Slug  string
	Price int
}

var products = []Product{
	{"t-shirt", "t-shirt", 80},
	{"cup", "cup", 20},
	{"book", "book", 50},
	{"pen", "pen", 10},
	{"powerbank", "powerbank", 200},
	{"hoody", "hoody", 300},
	{"umbrella", "umbrella", 200},
	{"socks", "socks", 10},
	{"wallet", "wallet", 50},
	{"pink-hoody", "pink-hoody", 500},
}

func main() {
	env, err := app.NewEnv()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	ctx := context.Background()

	conn, err := env.Database.PrimaryDB.Acquire(ctx)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		slog.Error(err.Error())
		return
	}

	// Сидируем товары
	if err := seedProducts(ctx, tx); err != nil {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil {
			err = errors.Join(err, rollbackErr)
		}
		slog.Error("Ошибка при сидировании", slog.Any("error", err))
		return
	}

	if err := tx.Commit(ctx); err != nil {
		slog.Error("Ошибка при коммите транзакции", slog.Any("error", err))
		return
	}
}
func seedProducts(ctx context.Context, tx pgx.Tx) error {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	for _, p := range products {
		query, args, err := psql.Insert("product").
			Columns("title", "slug", "price").
			Values(p.Title, p.Slug, p.Price).
			ToSql()
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, query, args...)
		if err != nil {
			return err
		}
	}
	return nil
}

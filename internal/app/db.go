package app

import (
	"context"

	"github.com/Sanchir01/avito-testovoe/pkg/db/connect"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	PrimaryDB *pgxpool.Pool
}

func NewDataBases(ctx context.Context, user, host, db, port string, maxAttempts int) (*Database, error) {
	pgxdb, err := connect.PGXNew(ctx, user, host, db, port, maxAttempts)
	if err != nil {
		return nil, err
	}

	return &Database{PrimaryDB: pgxdb}, nil
}

func (databases *Database) Close() error {
	databases.PrimaryDB.Close()
	return nil
}

package connect

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"time"
)

func PGXNew(ctx context.Context, user, host, db, port string, maxAttempts int) (*pgxpool.Pool, error) {

	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		user, os.Getenv("PASSWORD_POSTGRES"),
		host, port, db,
	)

	var pool *pgxpool.Pool

	err := DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		var err error
		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return err
		}
		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		return nil, err
	}

	return pool, nil
}
func DoWithTries(fn func() error, attemts int, delay time.Duration) (err error) {
	for attemts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attemts--
			continue
		}
		return nil
	}
	return
}

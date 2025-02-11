package user

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	repository *Repository
	primaryDB  *pgxpool.Pool
}

func NewService(r *Repository, db *pgxpool.Pool) *Service {
	return &Service{
		repository: r,
		primaryDB:  db,
	}
}

func (s *Service) Auth(ctx context.Context, email, password string) (string, error) {
	isExistUser, err := s.repository.GetUserByEmail(ctx, email)
	expirationTimeRefresh := time.Now().Add(14 * 24 * time.Hour)
	if err != nil {
		jwttoken, err := GenerateJwtToken(isExistUser.ID, expirationTimeRefresh)
		if err != nil {
			slog.Error("GenerateJwtToken err:", slog.Any("err", err))
			return "", err
		}
		return jwttoken, nil
	}
	conn, err := s.primaryDB.Acquire(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Release()
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = errors.Join(err, rollbackErr)
				return
			}
		}
	}()
	hashedPassword, err := GeneratePasswordHash(password)
	if err != nil {
		return "", err
	}
	id, err := s.repository.CreateUser(ctx, email, hashedPassword, tx)
	jwttoken, err := GenerateJwtToken(*id, expirationTimeRefresh)
	if err != nil {
		return "", err
	}
	return jwttoken, nil
}

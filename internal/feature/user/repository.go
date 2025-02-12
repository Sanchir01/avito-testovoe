package user

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	primaryDB *pgxpool.Pool
}

func NewRepository(primaryDB *pgxpool.Pool) *Repository {
	return &Repository{primaryDB: primaryDB}
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*DatabaseUser, error) {
	conn, err := r.primaryDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	query, arg, err := sq.
		Select("id,password,password").
		From("public.users").
		Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	var userDB DatabaseUser

	if err := conn.QueryRow(ctx, query, arg...).Scan(&userDB.ID, &userDB.Password, &userDB.Email); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("неправильный email")
		}
	}
	return &userDB, nil
}

func (r *Repository) CreateUser(ctx context.Context, email string, password []byte, tx pgx.Tx) (*uuid.UUID, error) {
	query, arg, err := sq.
		Insert("users").
		Columns("password", "email").
		Values(password, email).
		Suffix("RETURNING id,password").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}
	var (
		id           uuid.UUID
		passwordhash []byte
	)
	if err := tx.QueryRow(ctx, query, arg...).Scan(&id, &passwordhash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("неправильный email")
		}
	}
	return &id, nil
}

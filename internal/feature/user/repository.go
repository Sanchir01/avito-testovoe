package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/samber/lo"
	"log"

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
		Select("id, email,coin, version").
		From("public.users").
		Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}
	var userDB DatabaseUser
	if err := conn.QueryRow(ctx, query, arg...).Scan(&userDB.ID, &userDB.Email,
		&userDB.Coins, &userDB.Version,
	); err != nil {
		return nil, err
	}
	return &userDB, nil
}

func (r *Repository) CreateUser(ctx context.Context, email string, password []byte, tx pgx.Tx) (*uuid.UUID, error) {
	query, arg, err := sq.
		Insert("users").
		Columns("email", "password").
		Values(email, password).
		Suffix("RETURNING id, password").
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
		return nil, err
	}
	return &id, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*DatabaseUser, error) {
	conn, err := r.primaryDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query, arg, err := sq.
		Select("id, email,coin, version").
		From("public.users").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}
	var userDB DatabaseUser
	if err := conn.QueryRow(ctx, query, arg...).Scan(&userDB.ID, &userDB.Email,
		&userDB.Coins, &userDB.Version,
	); err != nil {
		return nil, err
	}
	return &userDB, nil
}

func (r *Repository) UpdateUserCoinByID(ctx context.Context, userID uuid.UUID, coins int64, tx pgx.Tx) error {
	query, arg, err := sq.
		Update("users").
		Set("coin", coins).
		Where(sq.Eq{"id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}
	_, exceErr := tx.Exec(ctx, query, arg...)

	if exceErr != nil {
		return fmt.Errorf("failed to update user coins: %s", exceErr.Error())
	}
	return nil
}
func (r *Repository) RecordPurchase(ctx context.Context, userID, productID uuid.UUID, tx pgx.Tx) error {
	query, args, err := sq.
		Insert("users_products").
		Columns("user_id", "product_id").
		Values(userID, productID).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return err
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to execute query: %v", err)
		return err
	}

	return nil
}
func (r *Repository) UpdateUserCoinByEmail(ctx context.Context, email string, coins int64, tx pgx.Tx) error {
	query, arg, err := sq.
		Update("users").
		Set("coin", coins).
		Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}
	_, exceErr := tx.Exec(ctx, query, arg...)

	if exceErr != nil {
		return fmt.Errorf("failed to update user coins: %s", exceErr.Error())
	}
	return nil
}

func (r *Repository) TransactionCoins(ctx context.Context, senderID, receiverID uuid.UUID, amount int, tx pgx.Tx) error {
	query, args, err := sq.
		Insert("transactions_coins").
		Columns("sender_id", "receiver_id", "amount").
		Values(senderID, receiverID, amount).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return err
	}
	return nil
}
func (r *Repository) GetAllProductBuyUsers(ctx context.Context, userID uuid.UUID) ([]Inventory, error) {
	conn, err := r.primaryDB.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()
	query, args, err := sq.
		Select("p.title", "COUNT(*)::int8 as quantity").
		From("product p").
		Join("users_products up ON p.id = up.product_id").
		Where(sq.Eq{"up.user_id": userID}).
		GroupBy("p.title").
		OrderBy("quantity DESC").
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
	var productsAndCounts []ProductCount
	for rows.Next() {
		var productAndCount ProductCount
		if err := rows.Scan(&productAndCount.Title, &productAndCount.Count); err != nil {
			return nil, err
		}
		productsAndCounts = append(productsAndCounts, productAndCount)
	}
	return lo.Map(productsAndCounts, func(productAndCount ProductCount, _ int) Inventory {
		return Inventory{
			Type:     productAndCount.Title,
			Quantity: productAndCount.Count,
		}
	}), nil
}

package user

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Sanchir01/avito-testovoe/internal/feature/product"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	repository        *Repository
	productRepository *product.Repository
	primaryDB         *pgxpool.Pool
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
	if err == nil {
		decodepass, err := base64.StdEncoding.DecodeString(isExistUser.Password)
		if err != nil {
			return "", err
		}
		verifypass := VerifyPassword(
			decodepass,
			password,
		)
		if verifypass {
			return "", fmt.Errorf("Неправильный пароль")
		}
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
	if err != nil {
		return "", err
	}
	jwttoken, err := GenerateJwtToken(*id, expirationTimeRefresh)
	if err != nil {

		return "", err
	}
	if err := tx.Commit(ctx); err != nil {
		return "", err
	}

	return jwttoken, nil
}

func (s *Service) BuyProduct(ctx context.Context, userID, productID uuid.UUID) error {
	slog.Error("attribute", slog.Any("userId", userID), slog.Any("productID", productID))

	slog.Error("product", product)
	user, err := s.repository.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	cantbuyproduct := user.Coins - product.Price

	if cantbuyproduct < 0 {
		return fmt.Errorf("Недостаточно монет")
	}
	conn, err := s.primaryDB.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {

		return err
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
	if err := s.repository.UpdateUserCoin(ctx, userID, cantbuyproduct, tx); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

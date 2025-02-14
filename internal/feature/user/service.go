package user

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/Sanchir01/avito-testovoe/pkg/lib/api"
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

func NewService(r *Repository, product *product.Repository, db *pgxpool.Pool) *Service {
	return &Service{
		repository:        r,
		primaryDB:         db,
		productRepository: product,
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
	productById, err := s.productRepository.GetProductByID(ctx, productID)
	user, err := s.repository.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	coins := user.Coins - productById.Price

	if coins < 0 {
		return api.ErrInsufficientCoins
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
	if err := s.repository.RecordPurchase(ctx, userID, productID, tx); err != nil {
		return err
	}
	if err := s.repository.UpdateUserCoinByID(ctx, userID, coins, tx); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Service) SendCoins(ctx context.Context, userId uuid.UUID, senderEmail string, amount int64) error {
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
	userSender, err := s.repository.GetUserByID(ctx, userId)
	if err != nil {
		return err
	}
	countpositive := IsPositiveCount(userSender.Coins, amount)

	if !countpositive {
		return fmt.Errorf("Недостаточно монет")
	}
	userReceiver, err := s.repository.GetUserByEmail(ctx, senderEmail)

	if err != nil {
		return err
	}

	if userSender.ID == userReceiver.ID {
		return api.ErrTransactionMyself
	}
	senderBalance := userSender.Coins - amount
	receiverBalance := userReceiver.Coins + amount
	slog.Info("balance start", userReceiver.Coins, "after plus", receiverBalance)
	if err := s.repository.UpdateUserCoinByID(ctx, userSender.ID, senderBalance, tx); err != nil {
		return err
	}
	if err := s.repository.UpdateUserCoinByID(ctx, userReceiver.ID, receiverBalance, tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetAllUserCoinsInfo(ctx context.Context, userId uuid.UUID) (*GetAllUserCoinsInfo, error) {
	user, err := s.repository.GetUserByID(ctx, userId)
	if err != nil {
		return nil, err
	}
	products, err := s.repository.GetAllProductBuyUsers(ctx, userId)
	if err != nil {
		return nil, err
	}
	slog.Info("get all user coin info", products)

	return &GetAllUserCoinsInfo{
		Coins:     user.Coins,
		Inventory: products,
	}, nil
}

func IsPositiveCount(a, b int64) bool {
	c := a + b
	if c < 0 {
		return false
	}
	return true
}

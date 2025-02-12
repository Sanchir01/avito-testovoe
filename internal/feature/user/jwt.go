package user

import (
	"errors"

	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	ID uuid.UUID `json:"id"`
	jwt.RegisteredClaims
}

func GenerateJwtToken(id uuid.UUID, expire time.Time) (string, error) {
	claim := &Claims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expire),
		},
	}

	secretKey := []byte(os.Getenv("JWT_SECRET"))
	tokens := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := tokens.SignedString(secretKey)

	if err != nil {
		slog.Error("GenerateJwtToken err:", slog.Any("err", err))
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")

}

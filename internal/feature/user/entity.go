package user

import (
	"github.com/google/uuid"
)

type DatabaseUser struct {
	ID       uuid.UUID `db:"id"`
	Email    string    `db:"email"`
	Password string    `db:"password"`
	Coins    int64     `db:"coin"`
	Version  int64     `db:"version"`
}

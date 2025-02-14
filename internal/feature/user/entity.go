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

type GetAllUserCoinsInfo struct {
	Coins        int64        `json:"coins"`
	Inventory    []Inventory  `json:"inventory"`
	CoinsHistory CoinsHistory `json:"coinsHistory"`
}
type Inventory struct {
	Type     string `json:"type"`
	Quantity int64  `json:"quantity"`
}
type CoinsHistory struct {
	Received SenderInfo `json:"received"`
	Send     SenderInfo `json:"send"`
}
type SenderInfo struct {
	Username string `json:"fromUser"`
	Amount   int64  `json:"amount"`
}
type ProductCount struct {
	Title string `db:"title"`
	Count int64  `db:"quantity"`
}

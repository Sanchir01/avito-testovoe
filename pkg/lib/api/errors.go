package api

import "errors"

var ErrInsufficientCoins = errors.New("недостаточно coin на балансе")

var ErrTransactionMyself = errors.New("нельзя отправлять деньги самому себе")

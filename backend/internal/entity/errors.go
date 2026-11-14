package entity

import "errors"

var (
	ErrOrderNotFound     = errors.New("order not found")
	ErrProductNotFound   = errors.New("product not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInsufficientCoins = errors.New("insufficient coins")
)

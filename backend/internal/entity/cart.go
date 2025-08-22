package entity

import (
	"time"

	"github.com/google/uuid"
)

type CartItem struct {
	ID        int       `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	ProductID int       `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type AddToCartResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type GetCartResponse struct {
	Items      []CartItem `json:"items"`
	TotalItems int        `json:"total_items"`
	Success    bool       `json:"success"`
	Message    string     `json:"message,omitempty"`
}

type RemoveFromCartResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type ClearCartResponse struct {
	Success      bool   `json:"success"`
	ItemsCleared int    `json:"items_cleared"`
	Message      string `json:"message,omitempty"`
}

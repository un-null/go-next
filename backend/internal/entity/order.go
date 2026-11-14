package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID             int         `json:"id"`
	UserID         uuid.UUID   `json:"user_id"`
	OrderNumber    string      `json:"order_number"`
	TotalAmount    float64     `json:"total_amount"`
	TotalCoinsUsed int         `json:"total_coins_used"`
	Status         OrderStatus `json:"status"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`

	Items []OrderItem `json:"items,omitempty"`
}

func (o *Order) IsCompleted() bool {
	return o.Status == OrderStatusCompleted
}

func (o *Order) IsPending() bool {
	return o.Status == OrderStatusPending
}

type OrderListResponse struct {
	Orders []Order `json:"orders"`
	Total  int     `json:"total"`
	Page   int     `json:"page"`
	Limit  int     `json:"limit"`
}

// CreateOrderRequest is what the client may send. The server decides the
// user, order number, amount, coins and status.
type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

// CreateOrderParams is built by the server and persisted as an order.
type CreateOrderParams struct {
	UserID         uuid.UUID
	OrderNumber    string
	TotalAmount    float64
	TotalCoinsUsed int
	Status         OrderStatus
}

type CreateOrderItemRequest struct {
	ProductID int `json:"product_id" validate:"required"`
	Quantity  int `json:"quantity" validate:"required,gt=0"`
}

type OrderQuery struct {
	UserID *uuid.UUID   `json:"user_id,omitempty"`
	Status *OrderStatus `json:"status,omitempty"`
	Page   int          `json:"page" validate:"min=1"`
	Limit  int          `json:"limit" validate:"min=1,max=100"`
}

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

type CreateOrderRequest struct {
	UserID         uuid.UUID                `json:"user_id" validate:"required"`
	OrderNumber    string                   `json:"order_number" validate:"required"`
	TotalAmount    float64                  `json:"total_amount" validate:"required"`
	TotalCoinsUsed int                      `json:"total_coins_used,omitempty"`
	Status         OrderStatus              `json:"status"`
	Items          []CreateOrderItemRequest `json:"items" validate:"required,min=1"`
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

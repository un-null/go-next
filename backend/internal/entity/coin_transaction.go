package entity

import (
	"time"

	"github.com/google/uuid"
)

type CoinTransaction struct {
	ID              int32     `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	TransactionType string    `json:"transaction_type"`
	Amount          int       `json:"amount"`
	BalanceAfter    int       `json:"balance_after"`
	OrderID         *int32    `json:"order_id,omitempty"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
}

func (ct *CoinTransaction) ToResponse() map[string]interface{} {
	response := map[string]interface{}{
		"id":               ct.ID,
		"user_id":          ct.UserID,
		"transaction_type": ct.TransactionType,
		"amount":           ct.Amount,
		"balance_after":    ct.BalanceAfter,
		"description":      ct.Description,
		"created_at":       ct.CreatedAt,
	}

	if ct.OrderID != nil {
		response["order_id"] = *ct.OrderID
	}

	return response
}

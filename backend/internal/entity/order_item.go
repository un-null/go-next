package entity

type OrderItem struct {
	ID           int     `json:"id"`
	OrderID      int     `json:"order_id"`
	ProductID    int     `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductPrice float64 `json:"product_price"`
	Quantity     int     `json:"quantity"`
	Subtotal     float64 `json:"subtotal"`

	Product *Product `json:"product,omitempty"`
}

// Helper methods
func (oi *OrderItem) CalculateSubtotal() float64 {
	return oi.ProductPrice * float64(oi.Quantity)
}

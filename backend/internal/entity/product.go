// entity/product.go
package entity

import "time"

type Product struct {
	ID            int       `json:"id"`
	CategoryID    int       `json:"category_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Price         float64   `json:"price"`
	StockQuantity int       `json:"stock_quantity"`
	ImageURL      string    `json:"image_url"`
	AverageRating float64   `json:"average_rating"`
	TotalComments int       `json:"total_comments"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// For joined queries (optional)
	Category *Category `json:"category,omitempty"`
}

// Helper method
func (p *Product) IsOutOfStock() bool {
	return p.StockQuantity == 0
}

type ProductListResponse struct {
	Products []Product `json:"products"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	Limit    int       `json:"limit"`
}

// For stock updates during purchase
type UpdateStockRequest struct {
	ProductID int `json:"product_id" validate:"required"`
	Quantity  int `json:"quantity" validate:"required,gt=0"` // Amount to reduce
}

// For search/filtering
type ProductQuery struct {
	CategoryID *int   `json:"category_id,omitempty"`
	Search     string `json:"search,omitempty"`
	Page       int    `json:"page" validate:"min=1"`
	Limit      int    `json:"limit" validate:"min=1,max=100"`
}

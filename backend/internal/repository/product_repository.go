package repository

import (
	"backend/internal/entity"
)

type ProductRepository interface {
	FindAllProduct() []entity.Product
}

type productRepository struct{}

func NewProductRepository() ProductRepository {
	return &productRepository{}
}

func (r *productRepository) FindAllProduct() []entity.Product {
	return []entity.Product{
		{ID: 1, Name: "Apple", Price: 1.2},
		{ID: 2, Name: "Banana", Price: 0.8},
	}
}

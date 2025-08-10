package repository

import (
	"backend/internal/entity"
)

type ProductRepository interface {
	GetAllProducts() []entity.Product
	GetProductById(id int) entity.Product
}

type productRepository struct {
	products []entity.Product
}

func NewProductRepository() ProductRepository {
	return &productRepository{
		products: []entity.Product{
			{ID: 1, Name: "Apple", Price: 1.2, Quantity: 3},
			{ID: 2, Name: "Banana", Price: 0.8, Quantity: 2},
		},
	}
}

func (r *productRepository) GetAllProducts() []entity.Product {
	return r.products
}

func (r *productRepository) GetProductById(id int) entity.Product {
	for _, p := range r.products {
		if p.ID == id {
			return p
		}
	}

	return entity.Product{}
}

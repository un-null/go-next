package usecase

import (
	"backend/internal/entity"
	"backend/internal/repository"
)

type ProductUseCase struct {
	repo repository.ProductRepository
}

func NewProductUseCase(r repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{repo: r}
}

func (u *ProductUseCase) GetAllProducts() []entity.Product {
	return u.repo.GetAllProducts()
}

func (u *ProductUseCase) GetProductById(id int) entity.Product {
	return u.repo.GetProductById(id)
}

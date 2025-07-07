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

func (u *ProductUseCase) ListProducts() []entity.Product {
	return u.repo.FindAllProduct()
}

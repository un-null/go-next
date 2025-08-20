package usecase

import (
	"context"

	"backend/internal/entity"
	"backend/internal/repository"
)

type ProductUseCase struct {
	repo repository.ProductRepository
}

func NewProductUseCase(r repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{repo: r}
}

func (u *ProductUseCase) GetAllProducts(ctx context.Context, page, limit int) ([]entity.Product, error) {
	return u.repo.GetAllProducts(ctx, page, limit)
}

func (u *ProductUseCase) GetProductByID(ctx context.Context, id int) (*entity.Product, error) {
	return u.repo.GetProductByID(ctx, id)
}

func (u *ProductUseCase) GetProductsByCategory(ctx context.Context, categoryID, page, limit int) ([]entity.Product, error) {
	return u.repo.GetProductsByCategory(ctx, categoryID, page, limit)
}

func (u *ProductUseCase) UpdateStock(ctx context.Context, productID, newStock int) error {
	return u.repo.UpdateStock(ctx, productID, newStock)
}

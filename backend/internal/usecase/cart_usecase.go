package usecase

import (
	"context"

	"backend/internal/entity"
	"backend/internal/repository"

	"github.com/google/uuid"
)

type CartUseCase struct {
	repo repository.CartRepository
}

func NewCartUseCase(r repository.CartRepository) *CartUseCase {
	return &CartUseCase{repo: r}
}

func (u *CartUseCase) AddToCart(ctx context.Context, userID uuid.UUID, productID int, quantity int) error {
	return u.repo.AddToCart(ctx, userID, productID, quantity)
}

func (u *CartUseCase) GetCartItems(ctx context.Context, userID uuid.UUID) ([]entity.CartItem, error) {
	return u.repo.GetCartItems(ctx, userID)
}

func (u *CartUseCase) RemoveFromCart(ctx context.Context, userID uuid.UUID, productID int) error {
	return u.repo.RemoveFromCart(ctx, userID, productID)
}

func (u *CartUseCase) ClearCart(ctx context.Context, userID uuid.UUID) error {
	return u.repo.ClearCart(ctx, userID)
}

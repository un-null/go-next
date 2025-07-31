package usecase

import (
	"backend/internal/entity"
	"backend/internal/repository"
)

type CartUseCase struct {
	repo repository.CartRepository
}

func NewCartUseCase(r repository.CartRepository) *CartUseCase {
	return &CartUseCase{repo: r}
}

func (u *CartUseCase) AddToCart(userID int, productID int, quantity int) error {
	return u.repo.AddToCart(userID, productID, quantity)
}

func (u *CartUseCase) GetCartItems(userID int) ([]entity.CartItem, error) {
	return u.repo.GetCartItems(userID)
}

func (u *CartUseCase) RemoveFromCart(userID int, productID int) error {
	return u.repo.RemoveFromCart(userID, productID)
}

func (u *CartUseCase) ClearCart(userID int) error {
	return u.repo.ClearCart(userID)
}

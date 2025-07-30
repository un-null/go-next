package repository

import (
	"backend/internal/entity"
	"errors"
)

type CartRepository interface {
	AddToCart(userID int, productID int, quantity int) error
	GetCartItems(userID int) ([]entity.CartItem, error)
	RemoveFromCart(userID int, productID int) error
	ClearCart(userID int) error
}

type cartRepository struct {
	cartData map[int][]entity.CartItem
}

func NewCartRepository() CartRepository {
	return &cartRepository{
		cartData: make(map[int][]entity.CartItem),
	}
}

func (r *cartRepository) AddToCart(userID int, productID int, quantity int) error {
	items := r.cartData[userID]

	for i, item := range items {
		if item.ProductID == productID {
			items[i].Quantity += quantity
			return nil
		}
	}

	newCartItem := entity.CartItem{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	}

	r.cartData[userID] = append(items, newCartItem)
	return nil
}

func (r *cartRepository) GetCartItems(userID int) ([]entity.CartItem, error) {
	return r.cartData[userID], nil
}

func (r *cartRepository) RemoveFromCart(userID int, productId int) error {
	items := r.cartData[userID]
	newItems := []entity.CartItem{}
	found := false

	for _, item := range items {
		if item.ProductID == productId {
			found = true
			continue
		}
		newItems = append(newItems, item)
	}

	if !found {
		return errors.New("item not found in cart")
	}

	r.cartData[userID] = newItems
	return nil
}

func (r *cartRepository) ClearCart(userId int) error {
	delete(r.cartData, userId)
	return nil
}

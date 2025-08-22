package repository

import (
	"backend/internal/database"
	"backend/internal/entity"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CartRepository interface {
	AddToCart(ctx context.Context, userID uuid.UUID, productID int, quantity int) error
	GetCartItems(ctx context.Context, userID uuid.UUID) ([]entity.CartItem, error)
	RemoveFromCart(ctx context.Context, userID uuid.UUID, productID int) error
	ClearCart(ctx context.Context, userID uuid.UUID) error
}

type cartRepository struct {
	queries *database.Queries
}

func NewCartRepository(queries *database.Queries) CartRepository {
	return &cartRepository{
		queries: queries,
	}
}

func (r *cartRepository) AddToCart(ctx context.Context, userID uuid.UUID, productID int, quantity int) error {
	// Check if item already exists in cart
	exists, err := r.queries.CheckCartItemExists(ctx, database.CheckCartItemExistsParams{
		UserID:    database.UUIDToPgtype(userID),
		ProductID: int32(productID),
	})
	if err != nil {
		return fmt.Errorf("failed to check if cart item exists: %w", err)
	}

	if exists {
		_, err = r.queries.UpdateCartItemQuantity(ctx, database.UpdateCartItemQuantityParams{
			UserID:    database.UUIDToPgtype(userID),
			ProductID: int32(productID),
			Quantity:  int32(quantity), // This should be the new total quantity, not addition
			UpdatedAt: database.TimeToPgtype(time.Now()),
		})
		if err != nil {
			return fmt.Errorf("failed to update cart item quantity: %w", err)
		}
		return nil
	}

	_, err = r.queries.CreateCartItem(ctx, database.CreateCartItemParams{
		UserID:    database.UUIDToPgtype(userID),
		ProductID: int32(productID),
		Quantity:  int32(quantity),
		CreatedAt: database.TimeToPgtype(time.Now()),
		UpdatedAt: database.TimeToPgtype(time.Now()),
	})
	if err != nil {
		return fmt.Errorf("failed to add item to cart: %w", err)
	}

	return nil
}

func (r *cartRepository) GetCartItems(ctx context.Context, userID uuid.UUID) ([]entity.CartItem, error) {
	dbCartItems, err := r.queries.GetCartItemsByUser(ctx, database.UUIDToPgtype(userID))

	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}

	cartItems := make([]entity.CartItem, len(dbCartItems))
	for i, dbItem := range dbCartItems {
		cartItems[i] = *dbCartItemToEntity(dbItem)
	}

	return cartItems, nil
}

func (r *cartRepository) RemoveFromCart(ctx context.Context, userID uuid.UUID, productID int) error {
	err := r.queries.DeleteCartItem(ctx, database.DeleteCartItemParams{
		UserID:    database.UUIDToPgtype(userID),
		ProductID: int32(productID),
	})

	if err != nil {
		return fmt.Errorf("failed to remove item from cart: %w", err)
	}

	return err
}

func (r *cartRepository) ClearCart(ctx context.Context, userID uuid.UUID) error {
	err := r.queries.DeleteAllCartItemsByUser(ctx, database.UUIDToPgtype(userID))

	if err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}

	return err
}

func dbCartItemToEntity(dbItem database.CartItem) *entity.CartItem {
	return &entity.CartItem{
		ID:        int(dbItem.ID),
		UserID:    database.PgtypeToUUID(dbItem.UserID),
		ProductID: int(dbItem.ProductID),
		Quantity:  int(dbItem.Quantity),
		CreatedAt: database.PgtypeToTime(dbItem.CreatedAt),
		UpdatedAt: database.PgtypeToTime(dbItem.UpdatedAt),
	}
}

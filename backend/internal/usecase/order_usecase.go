package usecase

import (
	"backend/internal/entity"
	"backend/internal/repository"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderUsecase interface {
	GetOrderById(ctx context.Context, id int) (*entity.Order, error)
	GetOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Order, error)
	CreateOrder(ctx context.Context, req entity.CreateOrderRequest) (*entity.Order, error)
}

type orderUsecase struct {
	orderRepo     repository.OrderRepository
	orderItemRepo repository.OrderItemRepository
	cartRepo      repository.CartRepository
	productRepo   repository.ProductRepository
	db            *pgxpool.Pool
}

func NewOrderUsecase(
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
	cartRepo repository.CartRepository,
	productRepo repository.ProductRepository,
	db *pgxpool.Pool,
) OrderUsecase {
	return &orderUsecase{
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
		cartRepo:      cartRepo,
		productRepo:   productRepo,
		db:            db,
	}
}

func (uc *orderUsecase) GetOrderById(ctx context.Context, id int) (*entity.Order, error) {
	return uc.orderRepo.GetOrderById(ctx, id)
}

func (uc *orderUsecase) GetOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Order, error) {
	return uc.orderRepo.GetOrdersByUserID(ctx, userID)
}

func (uc *orderUsecase) CreateOrder(ctx context.Context, req entity.CreateOrderRequest) (*entity.Order, error) {
	// Start transaction
	tx, err := uc.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	order, err := uc.orderRepo.CreateOrderWithTx(ctx, tx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	orderItems := make([]entity.OrderItem, 0, len(req.Items))
	for _, itemReq := range req.Items {
		// Get product details to populate order item
		product, err := uc.productRepo.GetProductByID(ctx, itemReq.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to get product %d: %w", itemReq.ProductID, err)
		}

		if product.StockQuantity < itemReq.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s: requested %d, available %d",
				product.Name, itemReq.Quantity, product.StockQuantity)
		}

		orderItem := entity.OrderItem{
			OrderID:      order.ID,
			ProductID:    itemReq.ProductID,
			ProductName:  product.Name,
			ProductPrice: product.Price,
			Quantity:     itemReq.Quantity,
		}

		orderItems = append(orderItems, orderItem)
	}

	_, err = uc.orderItemRepo.CreateOrderItemsWithTx(ctx, tx, orderItems)
	if err != nil {
		return nil, fmt.Errorf("failed to create order items: %w", err)
	}

	if err := uc.cartRepo.ClearCartWithTx(ctx, tx, req.UserID); err != nil {
		return nil, fmt.Errorf("failed to clear cart: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return order, nil
}

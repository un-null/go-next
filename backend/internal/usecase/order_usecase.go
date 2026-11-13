package usecase

import (
	"backend/internal/entity"
	"backend/internal/repository"
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderUsecase interface {
	GetOrderById(ctx context.Context, id int) (*entity.Order, error)
	GetOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Order, error)
	CreateOrder(ctx context.Context, userID uuid.UUID, req entity.CreateOrderRequest) (*entity.Order, error)
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

func (uc *orderUsecase) CreateOrder(ctx context.Context, userID uuid.UUID, req entity.CreateOrderRequest) (*entity.Order, error) {
	// Prices and totals are taken from the database, never from the client.
	orderItems := make([]entity.OrderItem, 0, len(req.Items))
	var totalCents int64
	for _, itemReq := range req.Items {
		product, err := uc.productRepo.GetProductByID(ctx, itemReq.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to get product %d: %w", itemReq.ProductID, err)
		}

		if product.StockQuantity < itemReq.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s: requested %d, available %d",
				product.Name, itemReq.Quantity, product.StockQuantity)
		}

		totalCents += int64(math.Round(product.Price*100)) * int64(itemReq.Quantity)
		orderItems = append(orderItems, entity.OrderItem{
			ProductID:    itemReq.ProductID,
			ProductName:  product.Name,
			ProductPrice: product.Price,
			Quantity:     itemReq.Quantity,
		})
	}

	tx, err := uc.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	order, err := uc.orderRepo.CreateOrderWithTx(ctx, tx, entity.CreateOrderParams{
		UserID:         userID,
		OrderNumber:    fmt.Sprintf("ORD-%d", time.Now().UnixNano()),
		TotalAmount:    float64(totalCents) / 100,
		TotalCoinsUsed: 0,
		Status:         entity.OrderStatusPending,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	for i := range orderItems {
		orderItems[i].OrderID = order.ID
	}

	if _, err := uc.orderItemRepo.CreateOrderItemsWithTx(ctx, tx, orderItems); err != nil {
		return nil, fmt.Errorf("failed to create order items: %w", err)
	}

	if err := uc.cartRepo.ClearCartWithTx(ctx, tx, userID); err != nil {
		return nil, fmt.Errorf("failed to clear cart: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	order.Items = orderItems
	return order, nil
}

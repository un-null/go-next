package usecase

import (
	"backend/internal/entity"
	"backend/internal/repository"
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderUsecase interface {
	GetOrderById(ctx context.Context, id int) (*entity.Order, error)
	GetOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Order, error)
	CreateOrder(ctx context.Context, userID uuid.UUID, req entity.CreateOrderRequest) (*entity.Order, error)
}

type stockReserver interface {
	ReserveStockWithTx(ctx context.Context, tx pgx.Tx, productID, quantity int) (entity.Product, error)
}

type userLocker interface {
	GetUserByIDForUpdateWithTx(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*entity.User, error)
	UpdateUserCoinsWithTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, coinsDelta int) (*entity.User, error)
}

type coinLedger interface {
	CreateCoinTransactionWithTx(ctx context.Context, tx pgx.Tx, params repository.CreateCoinTransactionParams) (*entity.CoinTransaction, error)
}

type orderUsecase struct {
	orderRepo     repository.OrderRepository
	orderItemRepo repository.OrderItemRepository
	cartRepo      repository.CartRepository
	productRepo   stockReserver
	userRepo      userLocker
	coinRepo      coinLedger
	db            *pgxpool.Pool
}

func NewOrderUsecase(
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
	cartRepo repository.CartRepository,
	productRepo stockReserver,
	userRepo userLocker,
	coinRepo coinLedger,
	db *pgxpool.Pool,
) OrderUsecase {
	return &orderUsecase{
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
		cartRepo:      cartRepo,
		productRepo:   productRepo,
		userRepo:      userRepo,
		coinRepo:      coinRepo,
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
	// Collapse duplicate lines so stock is checked once per product.
	quantities := make(map[int]int, len(req.Items))
	for _, item := range req.Items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("quantity must be positive for product %d", item.ProductID)
		}
		quantities[item.ProductID] += item.Quantity
	}
	if len(quantities) == 0 {
		return nil, fmt.Errorf("order must contain at least one item")
	}

	productIDs := make([]int, 0, len(quantities))
	for id := range quantities {
		productIDs = append(productIDs, id)
	}
	sort.Ints(productIDs)

	tx, err := uc.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Lock order is fixed: the user first, then products by ascending ID.
	user, err := uc.userRepo.GetUserByIDForUpdateWithTx(ctx, tx, userID)
	if err != nil {
		return nil, err
	}

	orderItems := make([]entity.OrderItem, 0, len(productIDs))
	var totalCents int64
	for _, productID := range productIDs {
		quantity := quantities[productID]

		product, err := uc.productRepo.ReserveStockWithTx(ctx, tx, productID, quantity)
		if err != nil {
			return nil, err
		}

		totalCents += int64(math.Round(product.Price*100)) * int64(quantity)
		orderItems = append(orderItems, entity.OrderItem{
			ProductID:    productID,
			ProductName:  product.Name,
			ProductPrice: product.Price,
			Quantity:     quantity,
		})
	}

	// 1 coin = 1 currency unit, rounded up.
	coins := int((totalCents + 99) / 100)
	if user.Coins < coins {
		return nil, fmt.Errorf("%w: need %d, have %d", entity.ErrInsufficientCoins, coins, user.Coins)
	}

	order, err := uc.orderRepo.CreateOrderWithTx(ctx, tx, entity.CreateOrderParams{
		UserID:         userID,
		OrderNumber:    fmt.Sprintf("ORD-%d", time.Now().UnixNano()),
		TotalAmount:    float64(totalCents) / 100,
		TotalCoinsUsed: coins,
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

	updatedUser, err := uc.userRepo.UpdateUserCoinsWithTx(ctx, tx, userID, -coins)
	if err != nil {
		return nil, err
	}

	orderID := int32(order.ID)
	if _, err := uc.coinRepo.CreateCoinTransactionWithTx(ctx, tx, repository.CreateCoinTransactionParams{
		UserID:          userID,
		TransactionType: "purchase",
		Amount:          -coins,
		BalanceAfter:    updatedUser.Coins,
		OrderID:         &orderID,
		Description:     "Order " + order.OrderNumber,
	}); err != nil {
		return nil, fmt.Errorf("failed to record coin transaction: %w", err)
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

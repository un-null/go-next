package repository

import (
	"backend/internal/database"
	"backend/internal/entity"
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderItemRepository interface {
	GetOrderItemById(ctx context.Context, id int) (*entity.OrderItem, error)
	GetOrderItemsByOrderID(ctx context.Context, orderID int) ([]entity.OrderItem, error)
	CreateOrderItem(ctx context.Context, req entity.OrderItem) (*entity.OrderItem, error)

	CreateOrderItemWithTx(ctx context.Context, tx pgx.Tx, req entity.OrderItem) (*entity.OrderItem, error)
	CreateOrderItemsWithTx(ctx context.Context, tx pgx.Tx, items []entity.OrderItem) ([]entity.OrderItem, error)
}

// OrderItemQuerier defines the minimal methods required from the database layer
// for order item repository operations
type OrderItemQuerier interface {
	GetOrderItemByID(ctx context.Context, id int32) (database.OrderItem, error)
	GetOrderItemsByOrderID(ctx context.Context, orderID int32) ([]database.OrderItem, error)
	CreateOrderItem(ctx context.Context, arg database.CreateOrderItemParams) (database.OrderItem, error)
	WithTx(tx pgx.Tx) *database.Queries
}

type orderItemRepository struct {
	queries OrderItemQuerier
	db      *pgxpool.Pool
}

func NewOrderItemRepository(queries OrderItemQuerier, db *pgxpool.Pool) OrderItemRepository {
	return &orderItemRepository{
		queries: queries,
		db:      db,
	}
}

func (r *orderItemRepository) GetOrderItemById(ctx context.Context, id int) (*entity.OrderItem, error) {
	dbItem, err := r.queries.GetOrderItemByID(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("order item not found")
		}
		return nil, err
	}

	item := &entity.OrderItem{
		ID:           int(dbItem.ID),
		OrderID:      int(dbItem.OrderID),
		ProductID:    int(dbItem.ProductID),
		ProductName:  dbItem.ProductName,
		ProductPrice: database.NumericToFloat64(dbItem.ProductPrice),
		Quantity:     int(dbItem.Quantity),
		Subtotal:     database.NumericToFloat64(dbItem.Subtotal),
	}

	return item, nil
}

func (r *orderItemRepository) GetOrderItemsByOrderID(ctx context.Context, orderID int) ([]entity.OrderItem, error) {
	dbItems, err := r.queries.GetOrderItemsByOrderID(ctx, int32(orderID))
	if err != nil {
		return nil, err
	}

	items := make([]entity.OrderItem, 0, len(dbItems))
	for _, dbItem := range dbItems {
		item := entity.OrderItem{
			ID:           int(dbItem.ID),
			OrderID:      int(dbItem.OrderID),
			ProductID:    int(dbItem.ProductID),
			ProductName:  dbItem.ProductName,
			ProductPrice: database.NumericToFloat64(dbItem.ProductPrice),
			Quantity:     int(dbItem.Quantity),
			Subtotal:     database.NumericToFloat64(dbItem.Subtotal),
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *orderItemRepository) CreateOrderItem(ctx context.Context, req entity.OrderItem) (*entity.OrderItem, error) {
	// Ensure subtotal is always calculated in Go before inserting
	subtotal := req.CalculateSubtotal()

	dbItem, err := r.queries.CreateOrderItem(ctx, database.CreateOrderItemParams{
		OrderID:      int32(req.OrderID),
		ProductID:    int32(req.ProductID),
		ProductName:  req.ProductName,
		ProductPrice: database.Float64ToNumeric(req.ProductPrice),
		Quantity:     int32(req.Quantity),
		Subtotal:     database.Float64ToNumeric(subtotal),
	})
	if err != nil {
		return nil, err
	}

	item := &entity.OrderItem{
		ID:           int(dbItem.ID),
		OrderID:      int(dbItem.OrderID),
		ProductID:    int(dbItem.ProductID),
		ProductName:  dbItem.ProductName,
		ProductPrice: database.NumericToFloat64(dbItem.ProductPrice),
		Quantity:     int(dbItem.Quantity),
		Subtotal:     database.NumericToFloat64(dbItem.Subtotal),
	}

	return item, nil
}

func (r *orderItemRepository) CreateOrderItemWithTx(ctx context.Context, tx pgx.Tx, req entity.OrderItem) (*entity.OrderItem, error) {
	// Create queries instance with transaction
	txQueries := r.queries.WithTx(tx)

	// Ensure subtotal is always calculated in Go before inserting
	subtotal := req.CalculateSubtotal()

	dbItem, err := txQueries.CreateOrderItem(ctx, database.CreateOrderItemParams{
		OrderID:      int32(req.OrderID),
		ProductID:    int32(req.ProductID),
		ProductName:  req.ProductName,
		ProductPrice: database.Float64ToNumeric(req.ProductPrice),
		Quantity:     int32(req.Quantity),
		Subtotal:     database.Float64ToNumeric(subtotal),
	})
	if err != nil {
		return nil, err
	}

	return r.dbOrderItemToEntity(dbItem), nil
}

func (r *orderItemRepository) CreateOrderItemsWithTx(ctx context.Context, tx pgx.Tx, items []entity.OrderItem) ([]entity.OrderItem, error) {
	if len(items) == 0 {
		return []entity.OrderItem{}, nil
	}

	createdItems := make([]entity.OrderItem, 0, len(items))

	for _, item := range items {
		createdItem, err := r.CreateOrderItemWithTx(ctx, tx, item)
		if err != nil {
			return nil, err
		}
		createdItems = append(createdItems, *createdItem)
	}

	return createdItems, nil
}

func (r *orderItemRepository) dbOrderItemToEntity(dbItem database.OrderItem) *entity.OrderItem {
	return &entity.OrderItem{
		ID:           int(dbItem.ID),
		OrderID:      int(dbItem.OrderID),
		ProductID:    int(dbItem.ProductID),
		ProductName:  dbItem.ProductName,
		ProductPrice: database.NumericToFloat64(dbItem.ProductPrice),
		Quantity:     int(dbItem.Quantity),
		Subtotal:     database.NumericToFloat64(dbItem.Subtotal),
	}
}

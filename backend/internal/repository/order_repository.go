package repository

import (
	"backend/internal/database"
	"backend/internal/entity"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository interface {
	GetOrderById(ctx context.Context, id int) (*entity.Order, error)
	GetOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Order, error)
	CreateOrder(ctx context.Context, req entity.CreateOrderRequest) (*entity.Order, error)
}

type orderRepository struct {
	queries *database.Queries
	db      *pgxpool.Pool
}

func NewOrderRepository(queries *database.Queries, db *pgxpool.Pool) OrderRepository {
	return &orderRepository{
		queries: queries,
		db:      db,
	}
}

func (r *orderRepository) GetOrderById(ctx context.Context, id int) (*entity.Order, error) {
	dbOrder, err := r.queries.GetOrderByID(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	order := &entity.Order{
		ID:             int(dbOrder.ID),
		UserID:         database.PgtypeToUUID(dbOrder.UserID),
		OrderNumber:    dbOrder.OrderNumber,
		TotalAmount:    database.NumericToFloat64(dbOrder.TotalAmount),
		TotalCoinsUsed: int(dbOrder.TotalCoinsUsed),
		Status:         entity.OrderStatus(database.NullOrderStatusToString(dbOrder.Status)),
		CreatedAt:      dbOrder.CreatedAt.Time,
		UpdatedAt:      dbOrder.UpdatedAt.Time,
	}

	return order, nil
}

func (r *orderRepository) GetOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Order, error) {
	dbOrders, err := r.queries.GetOrdersByUserID(ctx, database.UUIDToPgtype(userID))
	if err != nil {
		return nil, err
	}

	orders := make([]entity.Order, 0, len(dbOrders))
	for _, dbOrder := range dbOrders {
		order := entity.Order{
			ID:             int(dbOrder.ID),
			UserID:         database.PgtypeToUUID(dbOrder.UserID),
			OrderNumber:    dbOrder.OrderNumber,
			TotalAmount:    database.NumericToFloat64(dbOrder.TotalAmount),
			TotalCoinsUsed: int(dbOrder.TotalCoinsUsed),
			Status:         entity.OrderStatus(database.NullOrderStatusToString(dbOrder.Status)),
			CreatedAt:      dbOrder.CreatedAt.Time,
			UpdatedAt:      dbOrder.UpdatedAt.Time,
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *orderRepository) CreateOrder(ctx context.Context, req entity.CreateOrderRequest) (*entity.Order, error) {
	status := "success"

	dbOrder, err := r.queries.CreateOrder(ctx, database.CreateOrderParams{
		UserID:         database.UUIDToPgtype(req.UserID),
		OrderNumber:    req.OrderNumber,
		TotalAmount:    database.Float64ToNumeric(req.TotalAmount),
		TotalCoinsUsed: int32(req.TotalCoinsUsed),
		Status:         database.StringToNullOrderStatus(status),
	})
	if err != nil {
		return nil, err
	}

	order := &entity.Order{
		ID:             int(dbOrder.ID),
		UserID:         database.PgtypeToUUID(dbOrder.UserID),
		OrderNumber:    dbOrder.OrderNumber,
		TotalAmount:    database.NumericToFloat64(dbOrder.TotalAmount),
		TotalCoinsUsed: int(dbOrder.TotalCoinsUsed),
		Status:         entity.OrderStatus(database.NullOrderStatusToString(dbOrder.Status)),
		CreatedAt:      dbOrder.CreatedAt.Time,
		UpdatedAt:      dbOrder.UpdatedAt.Time,
	}

	return order, nil
}

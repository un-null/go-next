package repository

import (
	"backend/internal/database"
	"backend/internal/entity"
	"backend/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Helper function to create a mock database order
func createMockDBOrder(id int32, userID uuid.UUID, orderNumber string, totalAmount float64, totalCoins int, status string) database.Order {
	return database.Order{
		ID:             id,
		UserID:         pgtype.UUID{Bytes: userID, Valid: true},
		OrderNumber:    orderNumber,
		TotalAmount:    database.Float64ToNumeric(totalAmount),
		TotalCoinsUsed: int32(totalCoins),
		Status:         database.StringToNullOrderStatus(status),
		CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
}

// Setup test repository
func setupOrderTest(t *testing.T) (OrderRepository, *mocks.MockOrderQuerier) {
	mockQuerier := new(mocks.MockOrderQuerier)
	repo := NewOrderRepository(mockQuerier, nil)
	return repo, mockQuerier
}

func TestGetOrderById_Success(t *testing.T) {
	repo, mockQuerier := setupOrderTest(t)
	ctx := context.Background()

	userID := uuid.New()
	dbOrder := createMockDBOrder(1, userID, "ORD123", 100.5, 10, "PAID")

	mockQuerier.On("GetOrderByID", ctx, int32(1)).Return(dbOrder, nil)

	order, err := repo.GetOrderById(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, 1, order.ID)
	assert.Equal(t, userID, order.UserID)
	assert.Equal(t, "ORD123", order.OrderNumber)
	mockQuerier.AssertExpectations(t)
}

func TestGetOrderById_NotFound(t *testing.T) {
	repo, mockQuerier := setupOrderTest(t)
	ctx := context.Background()

	mockQuerier.On("GetOrderByID", ctx, int32(1)).Return(database.Order{}, errors.New("sql: no rows in result set"))

	order, err := repo.GetOrderById(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, order)
	mockQuerier.AssertExpectations(t)
}

func TestGetOrdersByUserID_Success(t *testing.T) {
	repo, mockQuerier := setupOrderTest(t)
	ctx := context.Background()

	userID := uuid.New()
	dbOrders := []database.Order{
		createMockDBOrder(1, userID, "ORD123", 50, 5, "PAID"),
		createMockDBOrder(2, userID, "ORD124", 75, 8, "PENDING"),
	}

	mockQuerier.On("GetOrdersByUserID", ctx, pgtype.UUID{Bytes: userID, Valid: true}).Return(dbOrders, nil)

	orders, err := repo.GetOrdersByUserID(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, orders, 2)
	assert.Equal(t, "ORD123", orders[0].OrderNumber)
	assert.Equal(t, "ORD124", orders[1].OrderNumber)
	mockQuerier.AssertExpectations(t)
}

func TestCreateOrder_Success(t *testing.T) {
	repo, mockQuerier := setupOrderTest(t)
	ctx := context.Background()

	userID := uuid.New()
	req := entity.CreateOrderRequest{
		UserID:         userID,
		OrderNumber:    "ORD125",
		TotalAmount:    200.5,
		TotalCoinsUsed: 20,
		Status:         entity.OrderStatus("PAID"),
	}

	dbOrder := createMockDBOrder(3, userID, "ORD125", 200.5, 20, "PAID")
	mockQuerier.On("CreateOrder", ctx, mock.Anything).Return(dbOrder, nil)

	order, err := repo.CreateOrder(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, 3, order.ID)
	assert.Equal(t, "ORD125", order.OrderNumber)
	mockQuerier.AssertExpectations(t)
}

func TestCreateOrderWithTx_Success(t *testing.T) {
	repo, mockQuerier := setupOrderTest(t)
	ctx := context.Background()

	userID := uuid.New()
	req := entity.CreateOrderRequest{
		UserID:         userID,
		OrderNumber:    "ORD126",
		TotalAmount:    300.5,
		TotalCoinsUsed: 30,
		Status:         entity.OrderStatus("PAID"),
	}

	dbOrder := createMockDBOrder(4, userID, "ORD126", 300.5, 30, "PAID")
	mockTx := new(pgx.Tx) // placeholder; will not actually execute SQL
	mockQuerier.On("CreateOrder", ctx, mock.Anything).Return(dbOrder, nil)

	order, err := repo.CreateOrderWithTx(ctx, *mockTx, req)

	assert.NoError(t, err)
	assert.Equal(t, 4, order.ID)
	assert.Equal(t, "ORD126", order.OrderNumber)
	mockQuerier.AssertExpectations(t)
}

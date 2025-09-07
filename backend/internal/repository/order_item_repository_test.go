package repository

import (
	"backend/internal/database"
	"backend/internal/entity"
	"backend/mocks"
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// helper to create a mock database order item
func createMockDBOrderItem(id, orderID, productID int32, name string, price float64, qty int32) database.OrderItem {
	return database.OrderItem{
		ID:           id,
		OrderID:      orderID,
		ProductID:    productID,
		ProductName:  name,
		ProductPrice: database.Float64ToNumeric(price),
		Quantity:     qty,
		Subtotal:     database.Float64ToNumeric(price * float64(qty)),
	}
}

func setupOrderItemTest(t *testing.T) (OrderItemRepository, *mocks.MockOrderItemQuerier) {
	q := new(mocks.MockOrderItemQuerier)
	repo := NewOrderItemRepository(q, nil)
	return repo, q
}

func TestGetOrderItemById_Success(t *testing.T) {
	repo, q := setupOrderItemTest(t)
	ctx := context.Background()

	dbItem := createMockDBOrderItem(1, 10, 100, "Apple", 120.0, 2)
	q.On("GetOrderItemByID", ctx, int32(1)).Return(dbItem, nil)

	item, err := repo.GetOrderItemById(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, 1, item.ID)
	assert.Equal(t, 10, item.OrderID)
	assert.Equal(t, "Apple", item.ProductName)
	q.AssertExpectations(t)
}

func TestGetOrderItemById_NotFound(t *testing.T) {
	repo, q := setupOrderItemTest(t)
	ctx := context.Background()

	q.On("GetOrderItemByID", ctx, int32(1)).Return(database.OrderItem{}, errors.New("sql: no rows in result set"))

	item, err := repo.GetOrderItemById(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, item)
	q.AssertExpectations(t)
}

func TestGetOrderItemsByOrderID_Success(t *testing.T) {
	repo, q := setupOrderItemTest(t)
	ctx := context.Background()

	dbItems := []database.OrderItem{
		createMockDBOrderItem(1, 10, 100, "Apple", 120.0, 2),
		createMockDBOrderItem(2, 10, 101, "Banana", 80.0, 3),
	}

	q.On("GetOrderItemsByOrderID", ctx, int32(10)).Return(dbItems, nil)

	items, err := repo.GetOrderItemsByOrderID(ctx, 10)

	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Equal(t, "Apple", items[0].ProductName)
	assert.Equal(t, "Banana", items[1].ProductName)
	q.AssertExpectations(t)
}

func TestCreateOrderItem_Success(t *testing.T) {
	repo, q := setupOrderItemTest(t)
	ctx := context.Background()

	req := entity.OrderItem{OrderID: 10, ProductID: 100, ProductName: "Apple", ProductPrice: 120.0, Quantity: 2}
	dbItem := createMockDBOrderItem(3, 10, 100, "Apple", 120.0, 2)
	q.On("CreateOrderItem", ctx, mock.Anything).Return(dbItem, nil)

	item, err := repo.CreateOrderItem(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, 3, item.ID)
	assert.Equal(t, 240.0, item.Subtotal)
	q.AssertExpectations(t)
}

func TestCreateOrderItemWithTx_Success(t *testing.T) {
	repo, q := setupOrderItemTest(t)
	ctx := context.Background()

	req := entity.OrderItem{OrderID: 10, ProductID: 100, ProductName: "Apple", ProductPrice: 120.0, Quantity: 2}
	dbItem := createMockDBOrderItem(4, 10, 100, "Apple", 120.0, 2)
	mockTx := new(pgx.Tx) // placeholder
	q.On("CreateOrderItem", ctx, mock.Anything).Return(dbItem, nil)

	item, err := repo.CreateOrderItemWithTx(ctx, *mockTx, req)

	assert.NoError(t, err)
	assert.Equal(t, 4, item.ID)
	assert.Equal(t, 240.0, item.Subtotal)
	q.AssertExpectations(t)
}

func TestCreateOrderItemsWithTx_Success(t *testing.T) {
	repo, q := setupOrderItemTest(t)
	ctx := context.Background()

	items := []entity.OrderItem{
		{OrderID: 10, ProductID: 100, ProductName: "Apple", ProductPrice: 120.0, Quantity: 2},
		{OrderID: 10, ProductID: 101, ProductName: "Banana", ProductPrice: 80.0, Quantity: 3},
	}

	q.On("CreateOrderItem", ctx, mock.Anything).Return(createMockDBOrderItem(5, 10, 100, "Apple", 120.0, 2), nil).Once()
	q.On("CreateOrderItem", ctx, mock.Anything).Return(createMockDBOrderItem(6, 10, 101, "Banana", 80.0, 3), nil).Once()

	mockTx := new(pgx.Tx)
	created, err := repo.CreateOrderItemsWithTx(ctx, *mockTx, items)

	assert.NoError(t, err)
	assert.Len(t, created, 2)
	assert.Equal(t, 5, created[0].ID)
	assert.Equal(t, 6, created[1].ID)
	q.AssertExpectations(t)
}

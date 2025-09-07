package usecase

import (
	"backend/internal/entity"
	"backend/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderItemUsecase_GetOrderItemById_Success(t *testing.T) {
	repo := new(mocks.MockOrderItemRepository)
	uc := NewOrderItemUsecase(repo)

	ctx := context.Background()
	expected := &entity.OrderItem{ID: 1, OrderID: 10, ProductID: 100, ProductName: "Apple", ProductPrice: 120, Quantity: 2}
	repo.On("GetOrderItemById", ctx, 1).Return(expected, nil)

	item, err := uc.GetOrderItemById(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, expected, item)
	repo.AssertExpectations(t)
}

func TestOrderItemUsecase_GetOrderItemById_Error(t *testing.T) {
	repo := new(mocks.MockOrderItemRepository)
	uc := NewOrderItemUsecase(repo)

	ctx := context.Background()
	expectedErr := errors.New("not found")
	repo.On("GetOrderItemById", ctx, 99).Return(nil, expectedErr)

	item, err := uc.GetOrderItemById(ctx, 99)

	assert.Error(t, err)
	assert.Nil(t, item)
	assert.Equal(t, expectedErr, err)
	repo.AssertExpectations(t)
}

func TestOrderItemUsecase_GetOrderItemsByOrderID_Success(t *testing.T) {
	repo := new(mocks.MockOrderItemRepository)
	uc := NewOrderItemUsecase(repo)

	ctx := context.Background()
	expected := []entity.OrderItem{{ID: 1, OrderID: 10}, {ID: 2, OrderID: 10}}
	repo.On("GetOrderItemsByOrderID", ctx, 10).Return(expected, nil)

	items, err := uc.GetOrderItemsByOrderID(ctx, 10)

	assert.NoError(t, err)
	assert.Equal(t, expected, items)
	repo.AssertExpectations(t)
}

func TestOrderItemUsecase_GetOrderItemsByOrderID_Error(t *testing.T) {
	repo := new(mocks.MockOrderItemRepository)
	uc := NewOrderItemUsecase(repo)

	ctx := context.Background()
	expectedErr := errors.New("db error")
	repo.On("GetOrderItemsByOrderID", ctx, 10).Return(nil, expectedErr)

	items, err := uc.GetOrderItemsByOrderID(ctx, 10)

	assert.Error(t, err)
	assert.Nil(t, items)
	assert.Equal(t, expectedErr, err)
	repo.AssertExpectations(t)
}

func TestOrderItemUsecase_CreateOrderItem_RecalculatesSubtotal(t *testing.T) {
	repo := new(mocks.MockOrderItemRepository)
	uc := NewOrderItemUsecase(repo)

	ctx := context.Background()
	input := entity.OrderItem{OrderID: 10, ProductID: 100, ProductName: "Apple", ProductPrice: 120, Quantity: 2}
	expectedOutput := input
	expectedOutput.Subtotal = expectedOutput.CalculateSubtotal()

	repo.On("CreateOrderItem", ctx, expectedOutput).Return(&expectedOutput, nil)

	created, err := uc.CreateOrderItem(ctx, input)

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput.Subtotal, created.Subtotal)
	repo.AssertExpectations(t)
}

func TestOrderItemUsecase_CreateOrderItem_Error(t *testing.T) {
	repo := new(mocks.MockOrderItemRepository)
	uc := NewOrderItemUsecase(repo)

	ctx := context.Background()
	input := entity.OrderItem{OrderID: 10, ProductID: 100, ProductName: "Apple", ProductPrice: 120, Quantity: 2}
	recalculated := input
	recalculated.Subtotal = recalculated.CalculateSubtotal()
	expectedErr := errors.New("insert error")

	repo.On("CreateOrderItem", ctx, recalculated).Return(nil, expectedErr)

	created, err := uc.CreateOrderItem(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, created)
	assert.Equal(t, expectedErr, err)
	repo.AssertExpectations(t)
}

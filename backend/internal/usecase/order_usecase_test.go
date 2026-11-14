package usecase

import (
	"backend/internal/entity"
	"backend/mocks"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestOrderUsecase_GetOrderById_Success(t *testing.T) {
	orderRepo := new(mocks.MockOrderRepository)
	orderItemRepo := new(mocks.MockOrderItemRepository)

	uc := NewOrderUsecase(orderRepo, orderItemRepo, nil, nil, nil, nil, nil)

	ctx := context.Background()
	expected := &entity.Order{ID: 1, OrderNumber: "ORD-1"}
	items := []entity.OrderItem{{ID: 10, OrderID: 1, ProductID: 100, Quantity: 2}}
	orderRepo.On("GetOrderById", ctx, 1).Return(expected, nil)
	orderItemRepo.On("GetOrderItemsByOrderID", ctx, 1).Return(items, nil)

	order, err := uc.GetOrderById(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, expected, order)
	assert.Equal(t, items, order.Items)
	orderRepo.AssertExpectations(t)
	orderItemRepo.AssertExpectations(t)
}

func TestOrderUsecase_GetOrderById_Error(t *testing.T) {
	orderRepo := new(mocks.MockOrderRepository)
	uc := NewOrderUsecase(orderRepo, nil, nil, nil, nil, nil, nil)

	ctx := context.Background()
	expectedErr := errors.New("not found")
	orderRepo.On("GetOrderById", ctx, 99).Return(nil, expectedErr)

	order, err := uc.GetOrderById(ctx, 99)

	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Equal(t, expectedErr, err)
	orderRepo.AssertExpectations(t)
}

func TestOrderUsecase_GetOrdersByUserID_Success(t *testing.T) {
	orderRepo := new(mocks.MockOrderRepository)
	uc := NewOrderUsecase(orderRepo, nil, nil, nil, nil, nil, nil)

	ctx := context.Background()
	userID := uuid.New()
	expected := []entity.Order{{ID: 1}, {ID: 2}}
	orderRepo.On("GetOrdersByUserID", ctx, userID).Return(expected, nil)

	orders, err := uc.GetOrdersByUserID(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expected, orders)
	orderRepo.AssertExpectations(t)
}

func TestOrderUsecase_GetOrdersByUserID_Error(t *testing.T) {
	orderRepo := new(mocks.MockOrderRepository)
	uc := NewOrderUsecase(orderRepo, nil, nil, nil, nil, nil, nil)

	ctx := context.Background()
	userID := uuid.New()
	expectedErr := errors.New("db error")
	orderRepo.On("GetOrdersByUserID", ctx, userID).Return(nil, expectedErr)

	orders, err := uc.GetOrdersByUserID(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, orders)
	assert.Equal(t, expectedErr, err)
	orderRepo.AssertExpectations(t)
}

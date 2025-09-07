package usecase

import (
	"backend/internal/entity"
	"backend/internal/repository"
	"context"
)

type OrderItemUsecase interface {
	GetOrderItemById(ctx context.Context, id int) (*entity.OrderItem, error)
	GetOrderItemsByOrderID(ctx context.Context, orderID int) ([]entity.OrderItem, error)
	CreateOrderItem(ctx context.Context, req entity.OrderItem) (*entity.OrderItem, error)
}

type orderItemUsecase struct {
	orderItemRepo repository.OrderItemRepository
}

func NewOrderItemUsecase(orderItemRepo repository.OrderItemRepository) OrderItemUsecase {
	return &orderItemUsecase{
		orderItemRepo: orderItemRepo,
	}
}

func (uc *orderItemUsecase) GetOrderItemById(ctx context.Context, id int) (*entity.OrderItem, error) {
	return uc.orderItemRepo.GetOrderItemById(ctx, id)
}

func (uc *orderItemUsecase) GetOrderItemsByOrderID(ctx context.Context, orderID int) ([]entity.OrderItem, error) {
	return uc.orderItemRepo.GetOrderItemsByOrderID(ctx, orderID)
}

func (uc *orderItemUsecase) CreateOrderItem(ctx context.Context, req entity.OrderItem) (*entity.OrderItem, error) {
	// enforce business logic: subtotal should always be recalculated here
	req.Subtotal = req.CalculateSubtotal()
	return uc.orderItemRepo.CreateOrderItem(ctx, req)
}

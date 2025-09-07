package usecase

import (
	"backend/internal/entity"
	"backend/internal/repository"
	"context"

	"github.com/google/uuid"
)

type OrderUsecase interface {
	GetOrderById(ctx context.Context, id int) (*entity.Order, error)
	GetOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Order, error)
	CreateOrder(ctx context.Context, req entity.CreateOrderRequest) (*entity.Order, error)
}

type orderUsecase struct {
	orderRepo repository.OrderRepository
}

func NewOrderUsecase(orderRepo repository.OrderRepository) OrderUsecase {
	return &orderUsecase{
		orderRepo: orderRepo,
	}
}

func (uc *orderUsecase) GetOrderById(ctx context.Context, id int) (*entity.Order, error) {
	return uc.orderRepo.GetOrderById(ctx, id)
}

func (uc *orderUsecase) GetOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Order, error) {
	return uc.orderRepo.GetOrdersByUserID(ctx, userID)
}

func (uc *orderUsecase) CreateOrder(ctx context.Context, req entity.CreateOrderRequest) (*entity.Order, error) {
	return uc.orderRepo.CreateOrder(ctx, req)
}

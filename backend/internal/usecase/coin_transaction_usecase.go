package usecase

import (
	"backend/internal/entity"
	"backend/internal/repository"
	"context"
	"errors"

	"github.com/google/uuid"
)

type CoinTransactionUseCase interface {
	ChargeUserCoins(ctx context.Context, userID uuid.UUID, amount int, description string, orderID *int32) (*entity.User, *entity.CoinTransaction, error)
	SpendUserCoins(ctx context.Context, userID uuid.UUID, amount int, description string, orderID *int32) (*entity.User, *entity.CoinTransaction, error)
	GetUserTransactions(ctx context.Context, userID uuid.UUID, page, limit int32) ([]*entity.CoinTransaction, error)
	GetTransactionByID(ctx context.Context, id int32) (*entity.CoinTransaction, error)
}

type coinTransactionUseCase struct {
	transactionRepo repository.CoinTransactionRepository
}

func NewCoinTransactionUseCase(transactionRepo repository.CoinTransactionRepository) CoinTransactionUseCase {
	return &coinTransactionUseCase{
		transactionRepo: transactionRepo,
	}
}

func (uc *coinTransactionUseCase) ChargeUserCoins(ctx context.Context, userID uuid.UUID, amount int, description string, orderID *int32) (*entity.User, *entity.CoinTransaction, error) {
	if amount <= 0 {
		return nil, nil, errors.New("amount must be positive")
	}

	if description == "" {
		return nil, nil, errors.New("description is required")
	}

	return uc.transactionRepo.ChargeUserCoins(ctx, userID, amount, description, orderID)
}

func (uc *coinTransactionUseCase) SpendUserCoins(ctx context.Context, userID uuid.UUID, amount int, description string, orderID *int32) (*entity.User, *entity.CoinTransaction, error) {
	if amount <= 0 {
		return nil, nil, errors.New("amount must be positive")
	}

	if description == "" {
		return nil, nil, errors.New("description is required")
	}

	return uc.transactionRepo.SpendUserCoins(ctx, userID, amount, description, orderID)
}

func (uc *coinTransactionUseCase) GetUserTransactions(ctx context.Context, userID uuid.UUID, page, limit int32) ([]*entity.CoinTransaction, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	return uc.transactionRepo.GetTransactionsByUserID(ctx, userID, limit, offset)
}

func (uc *coinTransactionUseCase) GetTransactionByID(ctx context.Context, id int32) (*entity.CoinTransaction, error) {
	if id <= 0 {
		return nil, errors.New("invalid transaction ID")
	}

	return uc.transactionRepo.GetTransactionByID(ctx, id)
}

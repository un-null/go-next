package repository

import (
	"backend/internal/database"
	"backend/internal/entity"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CoinTransactionRepository interface {
	CreateCoinTransaction(ctx context.Context, params CreateCoinTransactionParams) (*entity.CoinTransaction, error)
	GetTransactionsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*entity.CoinTransaction, error)
	GetTransactionByID(ctx context.Context, id int32) (*entity.CoinTransaction, error)
	ChargeUserCoins(ctx context.Context, userID uuid.UUID, amount int, description string, orderID *int32) (*entity.User, *entity.CoinTransaction, error)
	SpendUserCoins(ctx context.Context, userID uuid.UUID, amount int, description string, orderID *int32) (*entity.User, *entity.CoinTransaction, error)
}

type CreateCoinTransactionParams struct {
	UserID          uuid.UUID
	TransactionType string
	Amount          int
	BalanceAfter    int
	OrderID         *int32
	Description     string
}

type coinTransactionRepository struct {
	queries *database.Queries
	db      *pgxpool.Pool
}

func NewCoinTransactionRepository(queries *database.Queries, db *pgxpool.Pool) CoinTransactionRepository {
	return &coinTransactionRepository{
		queries: queries,
		db:      db,
	}
}

func (r *coinTransactionRepository) CreateCoinTransaction(ctx context.Context, params CreateCoinTransactionParams) (*entity.CoinTransaction, error) {
	var orderID pgtype.Int4
	if params.OrderID != nil {
		orderID = database.Int32ToPgtype(*params.OrderID)
	}

	dbTransaction, err := r.queries.CreateCoinTransaction(ctx, database.CreateCoinTransactionParams{
		UserID:          database.UUIDToPgtype(params.UserID),
		TransactionType: database.TransactionType(params.TransactionType),
		Amount:          int32(params.Amount),
		BalanceAfter:    int32(params.BalanceAfter),
		OrderID:         orderID,
		Description:     pgtype.Text{String: params.Description, Valid: params.Description != ""},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return dbTransactionToEntity(dbTransaction), nil
}

func (r *coinTransactionRepository) GetTransactionsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*entity.CoinTransaction, error) {
	dbTransactions, err := r.queries.GetCoinTransactionsByUserID(ctx, database.GetCoinTransactionsByUserIDParams{
		UserID: database.UUIDToPgtype(userID),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	transactions := make([]*entity.CoinTransaction, len(dbTransactions))
	for i, dbTx := range dbTransactions {
		transactions[i] = dbTransactionToEntity(dbTx)
	}

	return transactions, nil
}

func (r *coinTransactionRepository) GetTransactionByID(ctx context.Context, id int32) (*entity.CoinTransaction, error) {
	dbTransaction, err := r.queries.GetCoinTransactionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}

	return dbTransactionToEntity(dbTransaction), nil
}

// ChargeUserCoins - Add coins to user balance (like purchasing coins with real money)
func (r *coinTransactionRepository) ChargeUserCoins(ctx context.Context, userID uuid.UUID, amount int, description string, orderID *int32) (*entity.User, *entity.CoinTransaction, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	txQueries := r.queries.WithTx(tx)

	// Get current user
	user, err := txQueries.GetUserByID(ctx, database.UUIDToPgtype(userID))
	if err != nil {
		return nil, nil, fmt.Errorf("user not found: %w", err)
	}

	currentCoins := int(database.PgtypeToInt32(user.Coins))
	newBalance := currentCoins + amount // ADD coins for charging

	// Update user coins (positive amount)
	updatedUser, err := txQueries.UpdateUserCoins(ctx, database.UpdateUserCoinsParams{
		ID:    database.UUIDToPgtype(userID),
		Coins: database.Int32ToPgtype(int32(amount)), // Positive for addition
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update coins: %w", err)
	}

	// Create transaction record
	var orderIDPgtype pgtype.Int4
	if orderID != nil {
		orderIDPgtype = database.Int32ToPgtype(*orderID)
	}

	coinTx, err := txQueries.CreateCoinTransaction(ctx, database.CreateCoinTransactionParams{
		UserID:          database.UUIDToPgtype(userID),
		TransactionType: database.TransactionType("CHARGE"),
		Amount:          int32(amount), // Positive for charge
		BalanceAfter:    int32(newBalance),
		OrderID:         orderIDPgtype,
		Description:     pgtype.Text{String: description, Valid: description != ""},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Convert to entities
	userEntity := &entity.User{
		ID:        database.PgtypeToUUID(updatedUser.ID),
		Name:      updatedUser.Name,
		Email:     updatedUser.Email,
		Coins:     int(database.PgtypeToInt32(updatedUser.Coins)),
		CreatedAt: updatedUser.CreatedAt.Time,
		UpdatedAt: updatedUser.UpdatedAt.Time,
	}

	transactionEntity := &entity.CoinTransaction{
		ID:              coinTx.ID,
		UserID:          database.PgtypeToUUID(coinTx.UserID),
		TransactionType: string(coinTx.TransactionType),
		Amount:          int(coinTx.Amount),
		BalanceAfter:    int(coinTx.BalanceAfter),
		Description:     coinTx.Description.String,
		CreatedAt:       coinTx.CreatedAt.Time,
	}

	// Handle nullable OrderID
	if coinTx.OrderID.Valid {
		orderID := coinTx.OrderID.Int32
		transactionEntity.OrderID = &orderID
	}

	return userEntity, transactionEntity, nil
}

// SpendUserCoins - Deduct coins from user balance (like purchasing items with coins)
func (r *coinTransactionRepository) SpendUserCoins(ctx context.Context, userID uuid.UUID, amount int, description string, orderID *int32) (*entity.User, *entity.CoinTransaction, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	txQueries := r.queries.WithTx(tx)

	// Get current user
	user, err := txQueries.GetUserByID(ctx, database.UUIDToPgtype(userID))
	if err != nil {
		return nil, nil, fmt.Errorf("user not found: %w", err)
	}

	currentCoins := int(database.PgtypeToInt32(user.Coins))
	newBalance := currentCoins - amount // SUBTRACT coins for spending

	if newBalance < 0 {
		return nil, nil, fmt.Errorf("insufficient coins: have %d, need %d", currentCoins, amount)
	}

	// Update user coins (negative amount)
	updatedUser, err := txQueries.UpdateUserCoins(ctx, database.UpdateUserCoinsParams{
		ID:    database.UUIDToPgtype(userID),
		Coins: database.Int32ToPgtype(int32(-amount)), // Negative for deduction
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update coins: %w", err)
	}

	// Create transaction record
	var orderIDPgtype pgtype.Int4
	if orderID != nil {
		orderIDPgtype = database.Int32ToPgtype(*orderID)
	}

	coinTx, err := txQueries.CreateCoinTransaction(ctx, database.CreateCoinTransactionParams{
		UserID:          database.UUIDToPgtype(userID),
		TransactionType: database.TransactionType("SPEND"),
		Amount:          int32(-amount), // Negative for spend
		BalanceAfter:    int32(newBalance),
		OrderID:         orderIDPgtype,
		Description:     pgtype.Text{String: description, Valid: description != ""},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Convert to entities
	userEntity := &entity.User{
		ID:        database.PgtypeToUUID(updatedUser.ID),
		Name:      updatedUser.Name,
		Email:     updatedUser.Email,
		Coins:     int(database.PgtypeToInt32(updatedUser.Coins)),
		CreatedAt: updatedUser.CreatedAt.Time,
		UpdatedAt: updatedUser.UpdatedAt.Time,
	}

	transactionEntity := &entity.CoinTransaction{
		ID:              coinTx.ID,
		UserID:          database.PgtypeToUUID(coinTx.UserID),
		TransactionType: string(coinTx.TransactionType),
		Amount:          int(coinTx.Amount),
		BalanceAfter:    int(coinTx.BalanceAfter),
		Description:     coinTx.Description.String,
		CreatedAt:       coinTx.CreatedAt.Time,
	}

	// Handle nullable OrderID
	if coinTx.OrderID.Valid {
		orderID := coinTx.OrderID.Int32
		transactionEntity.OrderID = &orderID
	}

	return userEntity, transactionEntity, nil
}

func dbTransactionToEntity(dbTx database.CoinTransaction) *entity.CoinTransaction {
	transaction := &entity.CoinTransaction{
		ID:              dbTx.ID,
		UserID:          database.PgtypeToUUID(dbTx.UserID),
		TransactionType: string(dbTx.TransactionType),
		Amount:          int(dbTx.Amount),
		BalanceAfter:    int(dbTx.BalanceAfter),
		Description:     dbTx.Description.String,
		CreatedAt:       dbTx.CreatedAt.Time,
	}

	// Handle nullable OrderID
	if dbTx.OrderID.Valid {
		orderID := dbTx.OrderID.Int32
		transaction.OrderID = &orderID
	}

	return transaction
}

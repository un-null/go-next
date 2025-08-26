package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"backend/internal/database"
	"backend/internal/entity"
	"backend/mocks"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---- Test Repository ----

type testCoinTransactionRepository struct {
	queries *mocks.MockCoinTransactionQueries
	db      *mocks.MockDB
}

func NewTestCoinTransactionRepository(queries *mocks.MockCoinTransactionQueries, db *mocks.MockDB) CoinTransactionRepository {
	return &testCoinTransactionRepository{
		queries: queries,
		db:      db,
	}
}

func (r *testCoinTransactionRepository) CreateCoinTransaction(ctx context.Context, params CreateCoinTransactionParams) (*entity.CoinTransaction, error) {
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
		return nil, err
	}

	return dbTransactionToEntity(dbTransaction), nil
}

func (r *testCoinTransactionRepository) GetTransactionsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*entity.CoinTransaction, error) {
	dbTransactions, err := r.queries.GetCoinTransactionsByUserID(ctx, database.GetCoinTransactionsByUserIDParams{
		UserID: database.UUIDToPgtype(userID),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	transactions := make([]*entity.CoinTransaction, len(dbTransactions))
	for i, dbTx := range dbTransactions {
		transactions[i] = dbTransactionToEntity(dbTx)
	}

	return transactions, nil
}

func (r *testCoinTransactionRepository) GetTransactionByID(ctx context.Context, id int32) (*entity.CoinTransaction, error) {
	dbTransaction, err := r.queries.GetCoinTransactionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dbTransactionToEntity(dbTransaction), nil
}

func (r *testCoinTransactionRepository) ChargeUserCoins(ctx context.Context, userID uuid.UUID, amount int, description string, orderID *int32) (*entity.User, *entity.CoinTransaction, error) {
	// This would involve complex transaction mocking - simplified for testing
	return nil, nil, errors.New("not implemented in test")
}

func (r *testCoinTransactionRepository) SpendUserCoins(ctx context.Context, userID uuid.UUID, amount int, description string, orderID *int32) (*entity.User, *entity.CoinTransaction, error) {
	// This would involve complex transaction mocking - simplified for testing
	return nil, nil, errors.New("not implemented in test")
}

// ---- Helper Functions ----

func setupCoinTransactionTestRepository() (CoinTransactionRepository, *mocks.MockCoinTransactionQueries, *mocks.MockDB) {
	mockQueries := new(mocks.MockCoinTransactionQueries)
	mockDB := new(mocks.MockDB)
	repo := NewTestCoinTransactionRepository(mockQueries, mockDB)
	return repo, mockQueries, mockDB
}

func sampleDBCoinTransaction(id int32, userID uuid.UUID, txType string, amount int32) database.CoinTransaction {
	return database.CoinTransaction{
		ID:              id,
		UserID:          database.UUIDToPgtype(userID),
		TransactionType: database.TransactionType(txType),
		Amount:          amount,
		BalanceAfter:    1000 + amount, // Assume starting balance of 1000
		OrderID:         pgtype.Int4{Valid: false},
		Description:     pgtype.Text{String: "Test transaction", Valid: true},
		CreatedAt:       pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
}

func sampleUser(userID uuid.UUID, coins int32) database.User {
	return database.User{
		ID:        database.UUIDToPgtype(userID),
		Name:      "Test User",
		Email:     "test@example.com",
		Coins:     database.Int32ToPgtype(coins),
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
}

// ---- Tests ----

func TestCreateCoinTransaction_Success(t *testing.T) {
	repo, mockQ, _ := setupCoinTransactionTestRepository()

	userID := uuid.New()
	params := CreateCoinTransactionParams{
		UserID:          userID,
		TransactionType: "charge",
		Amount:          100,
		BalanceAfter:    1100,
		Description:     "Test charge",
	}

	expectedDB := sampleDBCoinTransaction(1, userID, "charge", 100)
	mockQ.On("CreateCoinTransaction", mock.Anything, mock.MatchedBy(func(p database.CreateCoinTransactionParams) bool {
		return database.PgtypeToUUID(p.UserID) == userID &&
			string(p.TransactionType) == "charge" &&
			p.Amount == 100 &&
			p.BalanceAfter == 1100
	})).Return(expectedDB, nil)

	transaction, err := repo.CreateCoinTransaction(context.Background(), params)

	assert.NoError(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, 1, transaction.ID)
	assert.Equal(t, userID, transaction.UserID)
	assert.Equal(t, "charge", transaction.TransactionType)
	assert.Equal(t, 100, transaction.Amount)
	assert.Equal(t, 1100, transaction.BalanceAfter)

	mockQ.AssertExpectations(t)
}

func TestCreateCoinTransaction_WithOrderID(t *testing.T) {
	repo, mockQ, _ := setupCoinTransactionTestRepository()

	userID := uuid.New()
	orderID := int32(123)
	params := CreateCoinTransactionParams{
		UserID:          userID,
		TransactionType: "purchase",
		Amount:          -50,
		BalanceAfter:    950,
		OrderID:         &orderID,
		Description:     "Purchase item",
	}

	expectedDB := sampleDBCoinTransaction(2, userID, "purchase", -50)
	expectedDB.OrderID = database.Int32ToPgtype(orderID)

	mockQ.On("CreateCoinTransaction", mock.Anything, mock.MatchedBy(func(p database.CreateCoinTransactionParams) bool {
		return p.OrderID.Valid && p.OrderID.Int32 == orderID
	})).Return(expectedDB, nil)

	transaction, err := repo.CreateCoinTransaction(context.Background(), params)

	assert.NoError(t, err)
	assert.NotNil(t, transaction)
	assert.NotNil(t, transaction.OrderID)
	assert.Equal(t, orderID, *transaction.OrderID)

	mockQ.AssertExpectations(t)
}

func TestCreateCoinTransaction_Error(t *testing.T) {
	repo, mockQ, _ := setupCoinTransactionTestRepository()

	userID := uuid.New()
	params := CreateCoinTransactionParams{
		UserID:          userID,
		TransactionType: "invalid_type",
		Amount:          100,
		BalanceAfter:    1100,
		Description:     "Test transaction",
	}

	mockQ.On("CreateCoinTransaction", mock.Anything, mock.Anything).
		Return(database.CoinTransaction{}, errors.New("invalid transaction type"))

	transaction, err := repo.CreateCoinTransaction(context.Background(), params)

	assert.Error(t, err)
	assert.Nil(t, transaction)
	assert.Contains(t, err.Error(), "invalid transaction type")

	mockQ.AssertExpectations(t)
}

func TestGetTransactionsByUserID_Success(t *testing.T) {
	repo, mockQ, _ := setupCoinTransactionTestRepository()

	userID := uuid.New()
	dbTransactions := []database.CoinTransaction{
		sampleDBCoinTransaction(1, userID, "charge", 100),
		sampleDBCoinTransaction(2, userID, "purchase", -50),
	}

	mockQ.On("GetCoinTransactionsByUserID", mock.Anything, database.GetCoinTransactionsByUserIDParams{
		UserID: database.UUIDToPgtype(userID),
		Limit:  10,
		Offset: 0,
	}).Return(dbTransactions, nil)

	transactions, err := repo.GetTransactionsByUserID(context.Background(), userID, 10, 0)

	assert.NoError(t, err)
	assert.Len(t, transactions, 2)
	assert.Equal(t, "charge", transactions[0].TransactionType)
	assert.Equal(t, "purchase", transactions[1].TransactionType)

	mockQ.AssertExpectations(t)
}

func TestGetTransactionsByUserID_WithPagination(t *testing.T) {
	repo, mockQ, _ := setupCoinTransactionTestRepository()

	userID := uuid.New()
	dbTransactions := []database.CoinTransaction{
		sampleDBCoinTransaction(3, userID, "refund", 25),
	}

	// Test page 2 with limit 5 (offset = 5)
	mockQ.On("GetCoinTransactionsByUserID", mock.Anything, database.GetCoinTransactionsByUserIDParams{
		UserID: database.UUIDToPgtype(userID),
		Limit:  5,
		Offset: 5,
	}).Return(dbTransactions, nil)

	transactions, err := repo.GetTransactionsByUserID(context.Background(), userID, 5, 5)

	assert.NoError(t, err)
	assert.Len(t, transactions, 1)
	assert.Equal(t, "refund", transactions[0].TransactionType)

	mockQ.AssertExpectations(t)
}

func TestGetTransactionsByUserID_Empty(t *testing.T) {
	repo, mockQ, _ := setupCoinTransactionTestRepository()

	userID := uuid.New()

	mockQ.On("GetCoinTransactionsByUserID", mock.Anything, database.GetCoinTransactionsByUserIDParams{
		UserID: database.UUIDToPgtype(userID),
		Limit:  10,
		Offset: 0,
	}).Return([]database.CoinTransaction{}, nil)

	transactions, err := repo.GetTransactionsByUserID(context.Background(), userID, 10, 0)

	assert.NoError(t, err)
	assert.Len(t, transactions, 0)

	mockQ.AssertExpectations(t)
}

func TestGetTransactionByID_Found(t *testing.T) {
	repo, mockQ, _ := setupCoinTransactionTestRepository()

	userID := uuid.New()
	expectedDB := sampleDBCoinTransaction(1, userID, "charge", 100)

	mockQ.On("GetCoinTransactionByID", mock.Anything, int32(1)).Return(expectedDB, nil)

	transaction, err := repo.GetTransactionByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, 1, transaction.ID)
	assert.Equal(t, userID, transaction.UserID)
	assert.Equal(t, "charge", transaction.TransactionType)

	mockQ.AssertExpectations(t)
}

func TestGetTransactionByID_NotFound(t *testing.T) {
	repo, mockQ, _ := setupCoinTransactionTestRepository()

	mockQ.On("GetCoinTransactionByID", mock.Anything, int32(999)).
		Return(database.CoinTransaction{}, errors.New("transaction not found"))

	transaction, err := repo.GetTransactionByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, transaction)
	assert.Contains(t, err.Error(), "transaction not found")

	mockQ.AssertExpectations(t)
}

func TestGetTransactionsByUserID_DatabaseError(t *testing.T) {
	repo, mockQ, _ := setupCoinTransactionTestRepository()

	userID := uuid.New()

	mockQ.On("GetCoinTransactionsByUserID", mock.Anything, mock.Anything).
		Return(nil, errors.New("database connection error"))

	transactions, err := repo.GetTransactionsByUserID(context.Background(), userID, 10, 0)

	assert.Error(t, err)
	assert.Nil(t, transactions)
	assert.Contains(t, err.Error(), "database connection error")

	mockQ.AssertExpectations(t)
}

// Test the conversion function
func TestDbTransactionToEntity(t *testing.T) {
	userID := uuid.New()
	orderID := int32(123)
	now := time.Now()

	dbTx := database.CoinTransaction{
		ID:              1,
		UserID:          database.UUIDToPgtype(userID),
		TransactionType: database.TransactionType("charge"),
		Amount:          100,
		BalanceAfter:    1100,
		OrderID:         database.Int32ToPgtype(orderID),
		Description:     pgtype.Text{String: "Test transaction", Valid: true},
		CreatedAt:       pgtype.Timestamptz{Time: now, Valid: true},
	}

	entity := dbTransactionToEntity(dbTx)

	assert.Equal(t, 1, entity.ID)
	assert.Equal(t, userID, entity.UserID)
	assert.Equal(t, "charge", entity.TransactionType)
	assert.Equal(t, 100, entity.Amount)
	assert.Equal(t, 1100, entity.BalanceAfter)
	assert.NotNil(t, entity.OrderID)
	assert.Equal(t, orderID, *entity.OrderID)
	assert.Equal(t, "Test transaction", entity.Description)
	assert.Equal(t, now, entity.CreatedAt)
}

func TestDbTransactionToEntity_NoOrderID(t *testing.T) {
	userID := uuid.New()
	now := time.Now()

	dbTx := database.CoinTransaction{
		ID:              2,
		UserID:          database.UUIDToPgtype(userID),
		TransactionType: database.TransactionType("purchase"),
		Amount:          -50,
		BalanceAfter:    950,
		OrderID:         pgtype.Int4{Valid: false}, // No order ID
		Description:     pgtype.Text{String: "Purchase", Valid: true},
		CreatedAt:       pgtype.Timestamptz{Time: now, Valid: true},
	}

	entity := dbTransactionToEntity(dbTx)

	assert.Equal(t, 2, entity.ID)
	assert.Equal(t, userID, entity.UserID)
	assert.Equal(t, "purchase", entity.TransactionType)
	assert.Equal(t, -50, entity.Amount)
	assert.Equal(t, 950, entity.BalanceAfter)
	assert.Nil(t, entity.OrderID) // Should be nil
	assert.Equal(t, "Purchase", entity.Description)
}

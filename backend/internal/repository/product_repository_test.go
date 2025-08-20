package repository

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"backend/internal/database"
	"backend/internal/entity"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---- Mock Queries Interface ----

type ProductQueriesInterface interface {
	ListProducts(ctx context.Context, arg database.ListProductsParams) ([]database.Product, error)
	GetProductByID(ctx context.Context, id int32) (database.Product, error)
	ListProductsByCategory(ctx context.Context, arg database.ListProductsByCategoryParams) ([]database.Product, error)
	UpdateProductStock(ctx context.Context, arg database.UpdateProductStockParams) (database.Product, error)
}

type MockQueries struct {
	mock.Mock
}

func (m *MockQueries) ListProducts(ctx context.Context, arg database.ListProductsParams) ([]database.Product, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]database.Product), args.Error(1)
}

func (m *MockQueries) GetProductByID(ctx context.Context, id int32) (database.Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(database.Product), args.Error(1)
}

func (m *MockQueries) ListProductsByCategory(ctx context.Context, arg database.ListProductsByCategoryParams) ([]database.Product, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]database.Product), args.Error(1)
}

func (m *MockQueries) UpdateProductStock(ctx context.Context, arg database.UpdateProductStockParams) (database.Product, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(database.Product), args.Error(1)
}

// ---- Test Repository (uses the same interface as main repo) ----

type testProductRepository struct {
	queries ProductQueriesInterface
}

func NewTestProductRepository(queries ProductQueriesInterface) ProductRepository {
	return &testProductRepository{
		queries: queries,
	}
}

func (r *testProductRepository) GetAllProducts(ctx context.Context, page, limit int) ([]entity.Product, error) {
	offset := (page - 1) * limit

	dbProducts, err := r.queries.ListProducts(ctx, database.ListProductsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	products := make([]entity.Product, len(dbProducts))
	for i, dbProduct := range dbProducts {
		products[i] = *dbProductToEntity(dbProduct) // Use existing function
	}

	return products, nil
}

func (r *testProductRepository) GetProductByID(ctx context.Context, id int) (*entity.Product, error) {
	dbProduct, err := r.queries.GetProductByID(ctx, int32(id))
	if err != nil {
		return nil, err
	}

	return dbProductToEntity(dbProduct), nil // Use existing function
}

func (r *testProductRepository) GetProductsByCategory(ctx context.Context, categoryID, page, limit int) ([]entity.Product, error) {
	offset := (page - 1) * limit

	dbProducts, err := r.queries.ListProductsByCategory(ctx, database.ListProductsByCategoryParams{
		CategoryID: int32(categoryID),
		Limit:      int32(limit),
		Offset:     int32(offset),
	})
	if err != nil {
		return nil, err
	}

	products := make([]entity.Product, len(dbProducts))
	for i, dbProduct := range dbProducts {
		products[i] = *dbProductToEntity(dbProduct) // Use existing function
	}

	return products, nil
}

func (r *testProductRepository) UpdateStock(ctx context.Context, productID, newStock int) error {
	_, err := r.queries.UpdateProductStock(ctx, database.UpdateProductStockParams{
		ID:            int32(productID),
		StockQuantity: int32(newStock),
	})
	return err
}

// ---- Helpers ----
func sampleDBProduct(id int32, name string) database.Product {
	return database.Product{
		ID:          id,
		CategoryID:  1,
		Name:        name,
		Description: pgtype.Text{String: "desc", Valid: true},
		Price: pgtype.Numeric{
			Int:   big.NewInt(1999), // $19.99
			Exp:   -2,               // 2 decimal places
			Valid: true,
		},
		StockQuantity: 10,
		ImageUrl:      pgtype.Text{String: "http://example.com", Valid: true},
		AverageRating: pgtype.Numeric{
			Int:   big.NewInt(450), // 4.5 rating
			Exp:   -2,
			Valid: true,
		},
		TotalComments: pgtype.Int4{Int32: 5, Valid: true},
		CreatedAt:     pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:     pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
}

// ---- Tests ----

func TestGetAllProducts(t *testing.T) {
	mockQ := new(MockQueries)
	repo := NewTestProductRepository(mockQ)

	dbProducts := []database.Product{sampleDBProduct(1, "Apple"), sampleDBProduct(2, "Banana")}

	mockQ.On("ListProducts", mock.Anything, database.ListProductsParams{Limit: 10, Offset: 0}).
		Return(dbProducts, nil)

	products, err := repo.GetAllProducts(context.Background(), 1, 10)
	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "Apple", products[0].Name)
	assert.Equal(t, "Banana", products[1].Name)
}

func TestGetProductByID_Found(t *testing.T) {
	mockQ := new(MockQueries)
	repo := NewTestProductRepository(mockQ)

	dbProduct := sampleDBProduct(1, "Apple")
	mockQ.On("GetProductByID", mock.Anything, int32(1)).Return(dbProduct, nil)

	product, err := repo.GetProductByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, "Apple", product.Name)
}

func TestGetProductByID_NotFound(t *testing.T) {
	mockQ := new(MockQueries)
	repo := NewTestProductRepository(mockQ)

	mockQ.On("GetProductByID", mock.Anything, int32(999)).Return(database.Product{}, errors.New("not found"))

	product, err := repo.GetProductByID(context.Background(), 999)
	assert.Error(t, err)
	assert.Nil(t, product)
}

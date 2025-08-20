package repository

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"backend/internal/database"
	"backend/internal/entity"
	"backend/mocks"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---- Test Repository (uses the same interface as main repo) ----

type testProductRepository struct {
	queries *mocks.MockProductQueries
}

func NewTestProductRepository(queries *mocks.MockProductQueries) ProductRepository {
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

// ---- Helper Functions ----

func setupProductTestRepository() (ProductRepository, *mocks.MockProductQueries) {
	mockQueries := new(mocks.MockProductQueries)
	repo := NewTestProductRepository(mockQueries)
	return repo, mockQueries
}

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
	repo, mockQ := setupProductTestRepository()

	dbProducts := []database.Product{sampleDBProduct(1, "Apple"), sampleDBProduct(2, "Banana")}

	mockQ.On("ListProducts", mock.Anything, database.ListProductsParams{Limit: 10, Offset: 0}).
		Return(dbProducts, nil)

	products, err := repo.GetAllProducts(context.Background(), 1, 10)
	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "Apple", products[0].Name)
	assert.Equal(t, "Banana", products[1].Name)

	mockQ.AssertExpectations(t)
}

func TestGetAllProducts_WithPagination(t *testing.T) {
	repo, mockQ := setupProductTestRepository()

	dbProducts := []database.Product{sampleDBProduct(3, "Orange")}

	// Test page 2 with limit 5
	mockQ.On("ListProducts", mock.Anything, database.ListProductsParams{Limit: 5, Offset: 5}).
		Return(dbProducts, nil)

	products, err := repo.GetAllProducts(context.Background(), 2, 5)
	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, "Orange", products[0].Name)

	mockQ.AssertExpectations(t)
}

func TestGetProductByID_Found(t *testing.T) {
	repo, mockQ := setupProductTestRepository()

	dbProduct := sampleDBProduct(1, "Apple")
	mockQ.On("GetProductByID", mock.Anything, int32(1)).Return(dbProduct, nil)

	product, err := repo.GetProductByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, "Apple", product.Name)
	assert.Equal(t, 1, product.ID)

	mockQ.AssertExpectations(t)
}

func TestGetProductByID_NotFound(t *testing.T) {
	repo, mockQ := setupProductTestRepository()

	mockQ.On("GetProductByID", mock.Anything, int32(999)).Return(database.Product{}, errors.New("not found"))

	product, err := repo.GetProductByID(context.Background(), 999)
	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Equal(t, "not found", err.Error())

	mockQ.AssertExpectations(t)
}

func TestGetProductsByCategory_Success(t *testing.T) {
	repo, mockQ := setupProductTestRepository()

	categoryID := 1
	dbProducts := []database.Product{
		sampleDBProduct(1, "Apple"),
		sampleDBProduct(2, "Banana"),
	}

	mockQ.On("ListProductsByCategory", mock.Anything, database.ListProductsByCategoryParams{
		CategoryID: int32(categoryID),
		Limit:      10,
		Offset:     0,
	}).Return(dbProducts, nil)

	products, err := repo.GetProductsByCategory(context.Background(), categoryID, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "Apple", products[0].Name)
	assert.Equal(t, "Banana", products[1].Name)

	mockQ.AssertExpectations(t)
}

func TestGetProductsByCategory_WithPagination(t *testing.T) {
	repo, mockQ := setupProductTestRepository()

	categoryID := 2
	dbProducts := []database.Product{sampleDBProduct(5, "Orange")}

	// Test page 3 with limit 2 (offset = 4)
	mockQ.On("ListProductsByCategory", mock.Anything, database.ListProductsByCategoryParams{
		CategoryID: int32(categoryID),
		Limit:      2,
		Offset:     4,
	}).Return(dbProducts, nil)

	products, err := repo.GetProductsByCategory(context.Background(), categoryID, 3, 2)
	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, "Orange", products[0].Name)

	mockQ.AssertExpectations(t)
}

func TestGetProductsByCategory_Empty(t *testing.T) {
	repo, mockQ := setupProductTestRepository()

	categoryID := 999
	mockQ.On("ListProductsByCategory", mock.Anything, database.ListProductsByCategoryParams{
		CategoryID: int32(categoryID),
		Limit:      10,
		Offset:     0,
	}).Return([]database.Product{}, nil)

	products, err := repo.GetProductsByCategory(context.Background(), categoryID, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, products, 0)

	mockQ.AssertExpectations(t)
}

func TestUpdateStock_Success(t *testing.T) {
	repo, mockQ := setupProductTestRepository()

	productID := 1
	newStock := 25
	updatedProduct := sampleDBProduct(int32(productID), "Apple")
	updatedProduct.StockQuantity = int32(newStock)

	mockQ.On("UpdateProductStock", mock.Anything, database.UpdateProductStockParams{
		ID:            int32(productID),
		StockQuantity: int32(newStock),
	}).Return(updatedProduct, nil)

	err := repo.UpdateStock(context.Background(), productID, newStock)
	assert.NoError(t, err)

	mockQ.AssertExpectations(t)
}

func TestUpdateStock_Error(t *testing.T) {
	repo, mockQ := setupProductTestRepository()

	productID := 999
	newStock := 25

	mockQ.On("UpdateProductStock", mock.Anything, database.UpdateProductStockParams{
		ID:            int32(productID),
		StockQuantity: int32(newStock),
	}).Return(database.Product{}, errors.New("product not found"))

	err := repo.UpdateStock(context.Background(), productID, newStock)
	assert.Error(t, err)
	assert.Equal(t, "product not found", err.Error())

	mockQ.AssertExpectations(t)
}

func TestUpdateStock_ZeroStock(t *testing.T) {
	repo, mockQ := setupProductTestRepository()

	productID := 1
	newStock := 0
	updatedProduct := sampleDBProduct(int32(productID), "Apple")
	updatedProduct.StockQuantity = 0

	mockQ.On("UpdateProductStock", mock.Anything, database.UpdateProductStockParams{
		ID:            int32(productID),
		StockQuantity: 0,
	}).Return(updatedProduct, nil)

	err := repo.UpdateStock(context.Background(), productID, newStock)
	assert.NoError(t, err)

	mockQ.AssertExpectations(t)
}

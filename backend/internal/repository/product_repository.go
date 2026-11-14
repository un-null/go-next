package repository

import (
	"context"
	"errors"
	"fmt"

	"backend/internal/database"
	"backend/internal/entity"

	"github.com/jackc/pgx/v5"
)

type ProductRepository interface {
	GetAllProducts(ctx context.Context, page, limit int) ([]entity.Product, error)
	GetProductByID(ctx context.Context, id int) (*entity.Product, error)
	GetProductsByCategory(ctx context.Context, categoryID, page, limit int) ([]entity.Product, error)
	UpdateStock(ctx context.Context, productID, newStock int) error
}

type productRepository struct {
	queries *database.Queries
}

func NewProductRepository(queries *database.Queries) *productRepository {
	return &productRepository{
		queries: queries,
	}
}

func (r *productRepository) GetAllProducts(ctx context.Context, page, limit int) ([]entity.Product, error) {
	offset := (page - 1) * limit

	dbProducts, err := r.queries.ListProducts(ctx, database.ListProductsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	products := make([]entity.Product, len(dbProducts))
	for i, dbProduct := range dbProducts {
		products[i] = *dbProductToEntity(dbProduct)
	}

	return products, nil
}

func (r *productRepository) GetProductByID(ctx context.Context, id int) (*entity.Product, error) {
	dbProduct, err := r.queries.GetProductByID(ctx, int32(id))
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	return dbProductToEntity(dbProduct), nil
}

func (r *productRepository) GetProductsByCategory(ctx context.Context, categoryID, page, limit int) ([]entity.Product, error) {
	offset := (page - 1) * limit

	dbProducts, err := r.queries.ListProductsByCategory(ctx, database.ListProductsByCategoryParams{
		CategoryID: int32(categoryID),
		Limit:      int32(limit),
		Offset:     int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get products by category: %w", err)
	}

	products := make([]entity.Product, len(dbProducts))
	for i, dbProduct := range dbProducts {
		products[i] = *dbProductToEntity(dbProduct)
	}

	return products, nil
}

func (r *productRepository) UpdateStock(ctx context.Context, productID, newStock int) error {
	_, err := r.queries.UpdateProductStock(ctx, database.UpdateProductStockParams{
		ID:            int32(productID),
		StockQuantity: int32(newStock),
	})
	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	return nil
}

func (r *productRepository) ReserveStockWithTx(ctx context.Context, tx pgx.Tx, productID, quantity int) (entity.Product, error) {
	txQueries := r.queries.WithTx(tx)

	dbProduct, err := txQueries.GetProductByIDForUpdate(ctx, int32(productID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Product{}, entity.ErrProductNotFound
		}
		return entity.Product{}, fmt.Errorf("failed to lock product: %w", err)
	}

	if int(dbProduct.StockQuantity) < quantity {
		return entity.Product{}, fmt.Errorf("%w: product %d has %d left", entity.ErrInsufficientStock, productID, dbProduct.StockQuantity)
	}

	dbProduct, err = txQueries.UpdateProductStock(ctx, database.UpdateProductStockParams{
		ID:            dbProduct.ID,
		StockQuantity: dbProduct.StockQuantity - int32(quantity),
	})
	if err != nil {
		return entity.Product{}, fmt.Errorf("failed to reserve stock: %w", err)
	}

	return *dbProductToEntity(dbProduct), nil
}

func dbProductToEntity(dbProduct database.Product) *entity.Product {
	return &entity.Product{
		ID:            int(dbProduct.ID),
		CategoryID:    int(dbProduct.CategoryID),
		Name:          dbProduct.Name,
		Description:   dbProduct.Description.String,
		Price:         database.NumericToFloat64(dbProduct.Price),
		StockQuantity: int(dbProduct.StockQuantity),
		ImageURL:      dbProduct.ImageUrl.String,
		AverageRating: database.NumericToFloat64(dbProduct.AverageRating),
		TotalComments: int(dbProduct.TotalComments.Int32),
		CreatedAt:     dbProduct.CreatedAt.Time,
		UpdatedAt:     dbProduct.UpdatedAt.Time,
	}
}

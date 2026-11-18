//go:build integration

package usecase

import (
	"context"
	"testing"

	"backend/internal/database"
	"backend/internal/repository"

	"github.com/stretchr/testify/require"
)

func TestListProductsPagesDoNotOverlap(t *testing.T) {
	pool := setupIntegrationDB(t)
	queries := database.New(pool)
	products := NewProductUseCase(repository.NewProductRepository(queries))

	ctx := context.Background()
	const limit = 7

	var ids []int
	for page := 1; page <= 3; page++ {
		found, err := products.GetAllProducts(ctx, "", page, limit)
		require.NoError(t, err)

		for _, product := range found {
			ids = append(ids, product.ID)
		}
	}

	want := []int{
		20, 19, 18, 17, 16, 15, 14, 13, 12, 11,
		10, 9, 8, 7, 6, 5, 4, 3, 2, 1,
	}
	require.Equal(t, want, ids)
}

func TestSearchProductsByKeyword(t *testing.T) {
	pool := setupIntegrationDB(t)
	queries := database.New(pool)
	products := NewProductUseCase(repository.NewProductRepository(queries))

	ctx := context.Background()

	found, err := products.GetAllProducts(ctx, "laptop", 1, 10)
	require.NoError(t, err)
	require.Len(t, found, 1)
	require.Equal(t, "Gaming Laptop G15", found[0].Name)

	found, err = products.GetAllProducts(ctx, "LAPTOP", 1, 10)
	require.NoError(t, err)
	require.Len(t, found, 1)

	found, err = products.GetAllProducts(ctx, "zzz", 1, 10)
	require.NoError(t, err)
	require.Empty(t, found)
}

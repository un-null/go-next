//go:build integration

package usecase

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"

	"backend/internal/database"
	"backend/internal/entity"
	"backend/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// These tests expect a throwaway PostgreSQL database. They refuse to run
// against a database whose name does not contain "test".
//
//	docker compose -f backend/compose.test.yml up -d
//	TEST_DATABASE_URL=postgres://postgres:postgres@localhost:5434/gonet_test?sslmode=disable \
//	  go test -tags=integration ./internal/usecase -run TestCheckout -v

func setupIntegrationDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	cfg, err := pgxpool.ParseConfig(dsn)
	require.NoError(t, err)
	require.Contains(t, strings.ToLower(cfg.ConnConfig.Database), "test",
		"refusing to run integration tests against a non-test database")

	ctx := context.Background()
	applyMigrations(t, ctx, dsn)

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	return pool
}

// applyMigrations rebuilds the schema from the same files used in production.
func applyMigrations(t *testing.T, ctx context.Context, dsn string) {
	t.Helper()

	cfg, err := pgx.ParseConfig(dsn)
	require.NoError(t, err)
	cfg.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	conn, err := pgx.ConnectConfig(ctx, cfg)
	require.NoError(t, err)
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, `DROP SCHEMA public CASCADE; CREATE SCHEMA public;`)
	require.NoError(t, err)

	files, err := filepath.Glob(filepath.Join("..", "..", "migrations", "*.up.sql"))
	require.NoError(t, err)
	sort.Strings(files)
	require.NotEmpty(t, files)

	for _, file := range files {
		migration, err := os.ReadFile(file)
		require.NoError(t, err)

		// Some files carry a goosedown section in the same file; never run it.
		sql := string(migration)
		if i := strings.Index(sql, "-- +goose Down"); i >= 0 {
			sql = sql[:i]
		}

		_, err = conn.Exec(ctx, sql)
		require.NoError(t, err, "migration %s", filepath.Base(file))
	}
}

func newCheckoutUsecase(pool *pgxpool.Pool) OrderUsecase {
	q := database.New(pool)

	return NewOrderUsecase(
		repository.NewOrderRepository(q, pool),
		repository.NewOrderItemRepository(q, pool),
		repository.NewCartRepository(q, pool),
		repository.NewProductRepository(q),
		repository.NewUserRepository(q, pool),
		repository.NewCoinTransactionRepository(q, pool),
		pool,
	)
}

type checkoutFixture struct {
	userID    uuid.UUID
	productID int
}

func seedCheckoutFixture(t *testing.T, pool *pgxpool.Pool, coins, stock int) checkoutFixture {
	t.Helper()

	ctx := context.Background()

	var userID uuid.UUID
	err := pool.QueryRow(ctx, `
		INSERT INTO users (name, email, password_hash, coins)
		VALUES ('Test User', $1, 'hash', $2)
		RETURNING id`,
		fmt.Sprintf("%s@example.com", uuid.NewString()), coins,
	).Scan(&userID)
	require.NoError(t, err)

	var productID int
	err = pool.QueryRow(ctx, `
		INSERT INTO products (category_id, name, description, price, stock_quantity, image_url)
		VALUES (1, 'Test Product', '', $1::numeric, $2, '')
		RETURNING id`,
		"19.99", stock,
	).Scan(&productID)
	require.NoError(t, err)

	return checkoutFixture{userID: userID, productID: productID}
}

func scalarInt(t *testing.T, pool *pgxpool.Pool, query string, args ...any) int {
	t.Helper()

	var value int
	require.NoError(t, pool.QueryRow(context.Background(), query, args...).Scan(&value))
	return value
}

func TestCheckout_Success(t *testing.T) {
	pool := setupIntegrationDB(t)
	ctx := context.Background()
	f := seedCheckoutFixture(t, pool, 5000, 10)

	_, err := pool.Exec(ctx,
		`INSERT INTO cart_items (user_id, product_id, quantity) VALUES ($1, $2, 3)`,
		f.userID, f.productID)
	require.NoError(t, err)

	uc := newCheckoutUsecase(pool)
	order, err := uc.CreateOrder(ctx, f.userID, entity.CreateOrderRequest{
		Items: []entity.CreateOrderItemRequest{
			{ProductID: f.productID, Quantity: 1},
			{ProductID: f.productID, Quantity: 2},
		},
	})
	require.NoError(t, err)

	assert.InDelta(t, 59.97, order.TotalAmount, 0.001)
	assert.Equal(t, 60, order.TotalCoinsUsed)
	assert.Equal(t, entity.OrderStatusPending, order.Status)
	require.Len(t, order.Items, 1)
	assert.Equal(t, 3, order.Items[0].Quantity)
	assert.InDelta(t, 59.97, order.Items[0].Subtotal, 0.001)

	assert.Equal(t, 7, scalarInt(t, pool, `SELECT stock_quantity FROM products WHERE id = $1`, f.productID))
	assert.Equal(t, 4940, scalarInt(t, pool, `SELECT coins FROM users WHERE id = $1`, f.userID))

	var amount, balanceAfter int
	require.NoError(t, pool.QueryRow(ctx,
		`SELECT amount, balance_after FROM coin_transactions WHERE user_id = $1`,
		f.userID).Scan(&amount, &balanceAfter))
	assert.Equal(t, -60, amount)
	assert.Equal(t, 4940, balanceAfter)

	assert.Equal(t, 1, scalarInt(t, pool, `SELECT count(*) FROM orders`))
	assert.Equal(t, 0, scalarInt(t, pool, `SELECT count(*) FROM cart_items WHERE user_id = $1`, f.userID))
}

func TestCheckout_InsufficientStock(t *testing.T) {
	pool := setupIntegrationDB(t)
	ctx := context.Background()
	f := seedCheckoutFixture(t, pool, 5000, 1)

	uc := newCheckoutUsecase(pool)
	order, err := uc.CreateOrder(ctx, f.userID, entity.CreateOrderRequest{
		Items: []entity.CreateOrderItemRequest{{ProductID: f.productID, Quantity: 2}},
	})

	require.Error(t, err)
	assert.Nil(t, order)
	assert.True(t, errors.Is(err, entity.ErrInsufficientStock), "got %v", err)

	assert.Equal(t, 1, scalarInt(t, pool, `SELECT stock_quantity FROM products WHERE id = $1`, f.productID))
	assert.Equal(t, 5000, scalarInt(t, pool, `SELECT coins FROM users WHERE id = $1`, f.userID))
	assert.Equal(t, 0, scalarInt(t, pool, `SELECT count(*) FROM orders`))
	assert.Equal(t, 0, scalarInt(t, pool, `SELECT count(*) FROM coin_transactions`))
}

func TestCheckout_InsufficientCoins(t *testing.T) {
	pool := setupIntegrationDB(t)
	ctx := context.Background()
	f := seedCheckoutFixture(t, pool, 10, 5)

	uc := newCheckoutUsecase(pool)
	order, err := uc.CreateOrder(ctx, f.userID, entity.CreateOrderRequest{
		Items: []entity.CreateOrderItemRequest{{ProductID: f.productID, Quantity: 1}},
	})

	require.Error(t, err)
	assert.Nil(t, order)
	assert.True(t, errors.Is(err, entity.ErrInsufficientCoins), "got %v", err)

	assert.Equal(t, 5, scalarInt(t, pool, `SELECT stock_quantity FROM products WHERE id = $1`, f.productID))
	assert.Equal(t, 10, scalarInt(t, pool, `SELECT coins FROM users WHERE id = $1`, f.userID))
	assert.Equal(t, 0, scalarInt(t, pool, `SELECT count(*) FROM orders`))
}

func TestCheckout_ConcurrentLastUnit(t *testing.T) {
	pool := setupIntegrationDB(t)
	ctx := context.Background()
	f := seedCheckoutFixture(t, pool, 1000, 1)

	uc := newCheckoutUsecase(pool)
	req := entity.CreateOrderRequest{
		Items: []entity.CreateOrderItemRequest{{ProductID: f.productID, Quantity: 1}},
	}

	start := make(chan struct{})
	errs := make([]error, 2)
	var wg sync.WaitGroup

	for i := range errs {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			_, errs[i] = uc.CreateOrder(ctx, f.userID, req)
		}(i)
	}

	close(start)
	wg.Wait()

	var created, soldOut int
	for _, err := range errs {
		switch {
		case err == nil:
			created++
		case errors.Is(err, entity.ErrInsufficientStock):
			soldOut++
		default:
			t.Fatalf("unexpected error: %v", err)
		}
	}

	assert.Equal(t, 1, created)
	assert.Equal(t, 1, soldOut)
	assert.Equal(t, 0, scalarInt(t, pool, `SELECT stock_quantity FROM products WHERE id = $1`, f.productID))
	assert.Equal(t, 1, scalarInt(t, pool, `SELECT count(*) FROM orders`))
	assert.Equal(t, 1, scalarInt(t, pool, `SELECT count(*) FROM coin_transactions`))
	assert.Equal(t, 980, scalarInt(t, pool, `SELECT coins FROM users WHERE id = $1`, f.userID))
}

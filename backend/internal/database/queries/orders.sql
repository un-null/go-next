-- name: CreateOrder :one
INSERT INTO "orders" (user_id, order_number, total_amount, total_coins_used, status)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, order_number, total_amount, total_coins_used, status, created_at, updated_at;

-- name: GetOrderByID :one
SELECT id, user_id, order_number, total_amount, total_coins_used, status, created_at, updated_at
FROM "orders"
WHERE id = $1;

-- name: GetOrdersByUserID :many
SELECT id, user_id, order_number, total_amount, total_coins_used, status, created_at, updated_at
FROM "orders"
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateOrderStatus :one
UPDATE "orders"
SET 
    status = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, user_id, order_number, total_amount, total_coins_used, status, created_at, updated_at;
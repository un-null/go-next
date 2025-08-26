-- name: CreateCoinTransaction :one
INSERT INTO coin_transactions (user_id, transaction_type, amount, balance_after, order_id, description)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, user_id, transaction_type, amount, balance_after, order_id, description, created_at;

-- name: GetCoinTransactionsByUserID :many
SELECT id, user_id, transaction_type, amount, balance_after, order_id, description, created_at
FROM coin_transactions
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetCoinTransactionByID :one
SELECT id, user_id, transaction_type, amount, balance_after, order_id, description, created_at
FROM coin_transactions
WHERE id = $1;
-- name: CreateUser :one
INSERT INTO users (name, email, password_hash, coins)
VALUES ($1, $2, $3, $4)
RETURNING id, name, email, coins, created_at, updated_at;
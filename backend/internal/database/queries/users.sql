-- queries/user.sql

-- name: CreateUser :one
INSERT INTO users (name, email, password_hash, coins)
VALUES ($1, $2, $3, $4)
RETURNING id, name, email, coins, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, name, email, password_hash, coins, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, name, email, password_hash, coins, created_at, updated_at
FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET 
    name = COALESCE(NULLIF($2, ''), name),
    email = COALESCE(NULLIF($3, ''), email),
    coins = COALESCE($4, coins),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, email, coins, created_at, updated_at;

-- name: UpdateUserPassword :exec
UPDATE users
SET 
    password_hash = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: UpdateUserCoins :one
UPDATE users
SET 
    coins = coins + $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, email, coins, created_at, updated_at;

-- name: CheckEmailExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);
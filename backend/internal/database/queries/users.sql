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

-- name: UpdateUserName :one
UPDATE users
SET 
    name = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, email, coins, created_at, updated_at;

-- name: UpdateUserEmail :one  
UPDATE users
SET 
    email = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, email, coins, created_at, updated_at;

-- name: UpdateUserCoins :one
UPDATE users
SET 
    coins = coins + $2,
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

-- name: CheckEmailExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);

-- name: CheckEmailExistsForOtherUser :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND id != $2);
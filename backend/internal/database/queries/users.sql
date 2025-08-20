-- name: CreateUser :one
INSERT INTO users (name, email, password_hash, coins)
VALUES ($1, $2, $3, COALESCE($4, 0))
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
    name = COALESCE($2, name),
    email = COALESCE($3, email),
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

-- name: ListUsers :many
SELECT id, name, email, coins, created_at, updated_at
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetUserCount :one
SELECT COUNT(*) FROM users;

-- name: UpdateUserCoins :one
UPDATE users
SET 
    coins = coins + $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, email, coins, created_at, updated_at;

-- name: CheckEmailExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);

-- Create indexes for better performance
CREATE INDEX idx_users_email ON users(email);

-- Create trigger to automatically update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
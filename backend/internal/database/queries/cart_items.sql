-- name: CreateCartItem :one
INSERT INTO cart_items (user_id, product_id, quantity, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, product_id, quantity, created_at, updated_at;

-- name: GetCartItemsByUser :many
SELECT id, user_id, product_id, quantity, created_at, updated_at
FROM cart_items
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateCartItemQuantity :one
UPDATE cart_items
SET 
    quantity = $3,
    updated_at = $4
WHERE user_id = $1 AND product_id = $2
RETURNING id, user_id, product_id, quantity, created_at, updated_at;

-- name: DeleteCartItem :exec
DELETE FROM cart_items
WHERE user_id = $1 AND product_id = $2;

-- name: DeleteAllCartItemsByUser :exec
DELETE FROM cart_items
WHERE user_id = $1;

-- name: CheckCartItemExists :one
SELECT EXISTS(SELECT 1 FROM cart_items WHERE user_id = $1 AND product_id = $2);
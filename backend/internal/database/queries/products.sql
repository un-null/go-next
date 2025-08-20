-- name: ListProducts :many
SELECT id, category_id, name, description, price, stock_quantity, image_url,
       average_rating, total_comments, created_at, updated_at
FROM products
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetProductByID :one
SELECT id, category_id, name, description, price, stock_quantity, image_url,
       average_rating, total_comments, created_at, updated_at
FROM products
WHERE id = $1;

-- name: ListProductsByCategory :many
SELECT id, category_id, name, description, price, stock_quantity, image_url,
       average_rating, total_comments, created_at, updated_at
FROM products
WHERE category_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateProductStock :one
UPDATE products
SET stock_quantity = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING id, category_id, name, description, price, stock_quantity, image_url,
          average_rating, total_comments, created_at, updated_at;
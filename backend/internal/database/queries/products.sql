-- name: ListProducts :many
SELECT id, category_id, name, description, price, stock_quantity, image_url,
       average_rating, total_comments, created_at, updated_at
FROM products
WHERE name ILIKE '%' || sqlc.arg('q')::text || '%'
ORDER BY created_at DESC, id DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: GetProductByID :one
SELECT id, category_id, name, description, price, stock_quantity, image_url,
       average_rating, total_comments, created_at, updated_at
FROM products
WHERE id = $1;

-- name: GetProductByIDForUpdate :one
SELECT id, category_id, name, description, price, stock_quantity, image_url,
       average_rating, total_comments, created_at, updated_at
FROM products
WHERE id = $1
FOR UPDATE;

-- name: ListProductsByCategory :many
SELECT id, category_id, name, description, price, stock_quantity, image_url,
       average_rating, total_comments, created_at, updated_at
FROM products
WHERE category_id = sqlc.arg('category_id')
  AND name ILIKE '%' || sqlc.arg('q')::text || '%'
ORDER BY created_at DESC, id DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: UpdateProductStock :one
UPDATE products
SET stock_quantity = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING id, category_id, name, description, price, stock_quantity, image_url,
          average_rating, total_comments, created_at, updated_at;
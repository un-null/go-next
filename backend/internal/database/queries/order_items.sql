-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, product_name, product_price, quantity, subtotal)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, order_id, product_id, product_name, product_price, quantity, subtotal;

-- name: GetOrderItemsByOrderID :many
SELECT id, order_id, product_id, product_name, product_price, quantity, subtotal
FROM order_items
WHERE order_id = $1;

-- name: GetOrderItemByID :one
SELECT id, order_id, product_id, product_name, product_price, quantity, subtotal
FROM order_items
WHERE id = $1;
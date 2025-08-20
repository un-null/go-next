-- name: GetAllCategories :many
SELECT id, name
FROM categories
ORDER BY name;

-- name: GetCategoryByID :one
SELECT id, name
FROM categories
WHERE id = $1;
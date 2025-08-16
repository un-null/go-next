-- Drop trigger first
DROP TRIGGER IF EXISTS update_products_updated_at ON products;

-- Drop indexes
DROP INDEX IF EXISTS idx_products_category_active;
DROP INDEX IF EXISTS idx_products_active;
DROP INDEX IF EXISTS idx_products_price;
DROP INDEX IF EXISTS idx_products_rating;

-- Drop table
DROP TABLE IF EXISTS products;
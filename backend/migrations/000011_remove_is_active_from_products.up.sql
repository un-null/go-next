-- Remove is_active column and related indexes
DROP INDEX IF EXISTS idx_products_category_active;
DROP INDEX IF EXISTS idx_products_active;

-- Remove the column
ALTER TABLE products DROP COLUMN is_active;

-- Recreate category index without is_active
CREATE INDEX idx_products_category ON products(category_id);
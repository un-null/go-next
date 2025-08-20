-- Add back is_active column
ALTER TABLE products ADD COLUMN is_active BOOLEAN DEFAULT TRUE;

-- Recreate original indexes
CREATE INDEX idx_products_category_active ON products(category_id, is_active);
CREATE INDEX idx_products_active ON products(is_active);
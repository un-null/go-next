-- Remove unnecessary columns from categories
ALTER TABLE categories DROP COLUMN IF EXISTS description;
ALTER TABLE categories DROP COLUMN IF EXISTS image_url;
ALTER TABLE categories DROP COLUMN IF EXISTS is_active;
ALTER TABLE categories DROP COLUMN IF EXISTS sort_order;
ALTER TABLE categories DROP COLUMN IF EXISTS created_at;
ALTER TABLE categories DROP COLUMN IF EXISTS updated_at;

-- Drop triggers that are no longer needed
DROP TRIGGER IF EXISTS update_categories_updated_at ON categories;

-- Drop indexes that are no longer needed
DROP INDEX IF EXISTS idx_categories_active_sort;
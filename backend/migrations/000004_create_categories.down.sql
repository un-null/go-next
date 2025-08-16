-- Drop trigger first
DROP TRIGGER IF EXISTS update_categories_updated_at ON categories;

-- Drop indexes
DROP INDEX IF EXISTS idx_categories_active_sort;

-- Drop table
DROP TABLE IF EXISTS categories;
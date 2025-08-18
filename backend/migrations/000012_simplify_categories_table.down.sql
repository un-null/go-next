-- Add back the columns if we need to rollback
ALTER TABLE categories 
ADD COLUMN description TEXT,
ADD COLUMN image_url VARCHAR(500),
ADD COLUMN is_active BOOLEAN DEFAULT TRUE,
ADD COLUMN sort_order INTEGER DEFAULT 0,
ADD COLUMN created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;

-- Recreate trigger
CREATE TRIGGER update_categories_updated_at 
    BEFORE UPDATE ON categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Recreate index
CREATE INDEX idx_categories_active_sort ON categories(is_active, sort_order);
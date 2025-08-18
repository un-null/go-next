-- Drop triggers first
DROP TRIGGER IF EXISTS update_product_rating_on_comment ON comments;
DROP TRIGGER IF EXISTS update_comments_updated_at ON comments;

-- Drop function
DROP FUNCTION IF EXISTS update_product_rating();

-- Drop indexes
DROP INDEX IF EXISTS idx_comments_product_id;
DROP INDEX IF EXISTS idx_comments_user_id;
DROP INDEX IF EXISTS idx_comments_rating;

-- Drop table
DROP TABLE IF EXISTS comments;
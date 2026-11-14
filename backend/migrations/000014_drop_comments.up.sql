-- Drop the unused comments feature.
-- Nothing in the API reads or writes comments; product ratings keep the
-- values seeded in 000013.
DROP TRIGGER IF EXISTS update_product_rating_on_comment ON comments;
DROP TRIGGER IF EXISTS update_comments_updated_at ON comments;
DROP FUNCTION IF EXISTS update_product_rating();
DROP TABLE IF EXISTS comments;

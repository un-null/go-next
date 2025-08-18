-- Drop trigger first
DROP TRIGGER IF EXISTS update_orders_updated_at ON orders;

-- Drop indexes
DROP INDEX IF EXISTS idx_orders_user_id;
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_created_at;

-- Drop table
DROP TABLE IF EXISTS orders;

-- Drop ENUM type
DROP TYPE IF EXISTS order_status;
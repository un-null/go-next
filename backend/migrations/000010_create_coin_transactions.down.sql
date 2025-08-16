-- Drop indexes
DROP INDEX IF EXISTS idx_coin_transactions_user_id;
DROP INDEX IF EXISTS idx_coin_transactions_type;
DROP INDEX IF EXISTS idx_coin_transactions_created_at;

-- Drop table
DROP TABLE IF EXISTS coin_transactions;

-- Drop ENUM type
DROP TYPE IF EXISTS transaction_type;
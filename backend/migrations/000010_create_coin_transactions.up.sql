-- Create ENUM type
CREATE TYPE transaction_type AS ENUM ('charge', 'purchase', 'refund');

-- Create coin_transactions table
CREATE TABLE coin_transactions (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    transaction_type transaction_type NOT NULL,
    amount INTEGER NOT NULL, -- Positive for charge/refund, negative for purchase
    balance_after INTEGER NOT NULL, -- User's coin balance after this transaction
    order_id INTEGER NULL REFERENCES orders(id),
    description VARCHAR(255), -- e.g., "Purchased Order #ORD-001"
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_coin_transactions_user_id ON coin_transactions(user_id);
CREATE INDEX idx_coin_transactions_type ON coin_transactions(transaction_type);
CREATE INDEX idx_coin_transactions_created_at ON coin_transactions(created_at);
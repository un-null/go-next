-- Create order_items table
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL REFERENCES products(id),
    product_name VARCHAR(255) NOT NULL, -- Snapshot at time of purchase
    product_price DECIMAL(10,2) NOT NULL, -- Snapshot at time of purchase
    quantity INTEGER NOT NULL,
    subtotal DECIMAL(10,2) NOT NULL -- quantity * product_price
);

-- Create indexes
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
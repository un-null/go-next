-- Create comments table
CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    rating INTEGER CHECK (rating >= 1 AND rating <= 5),
    comment TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_comments_product_id ON comments(product_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_rating ON comments(rating);

-- Create trigger for updated_at
CREATE TRIGGER update_comments_updated_at 
    BEFORE UPDATE ON comments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to update product rating when comments change
CREATE OR REPLACE FUNCTION update_product_rating()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE products 
    SET 
        average_rating = COALESCE((
            SELECT ROUND(AVG(rating::DECIMAL), 2)
            FROM comments 
            WHERE product_id = COALESCE(NEW.product_id, OLD.product_id)
        ), 0),
        total_comments = (
            SELECT COUNT(*)
            FROM comments 
            WHERE product_id = COALESCE(NEW.product_id, OLD.product_id)
        )
    WHERE id = COALESCE(NEW.product_id, OLD.product_id);
    
    RETURN COALESCE(NEW, OLD);
END;
$$ language 'plpgsql';

-- Create trigger for automatic rating updates
CREATE TRIGGER update_product_rating_on_comment
    AFTER INSERT OR UPDATE OR DELETE ON comments
    FOR EACH ROW EXECUTE FUNCTION update_product_rating();
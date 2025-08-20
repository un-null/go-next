-- +goose Up
-- Insert categories
INSERT INTO categories (id, name) VALUES
(1, 'Electronics'),
(2, 'Books'),
(3, 'Clothing');

-- Insert products
INSERT INTO products (category_id, name, description, price, stock_quantity, image_url, average_rating, total_comments, created_at, updated_at) VALUES
-- Electronics
(1, 'Smartphone X200', 'High-performance smartphone with OLED display', 699.99, 50, 'https://images.unsplash.com/photo-1510557880182-3d4d3cba35a5', 4.5, 120, NOW(), NOW()),
(1, 'Wireless Earbuds Pro', 'Noise-cancelling true wireless earbuds', 149.99, 200, 'https://images.unsplash.com/photo-1585386959984-a41552231693', 4.3, 85, NOW(), NOW()),
(1, 'Gaming Laptop G15', 'Powerful gaming laptop with RTX graphics', 1299.99, 30, 'https://images.unsplash.com/photo-1587202372775-98973f1b64b8', 4.7, 210, NOW(), NOW()),
(1, '4K Monitor UltraSharp', '27-inch 4K UHD IPS monitor', 399.99, 80, 'https://images.unsplash.com/photo-1517336714731-489689fd1ca8', 4.4, 60, NOW(), NOW()),
(1, 'Mechanical Keyboard RGB', 'Mechanical keyboard with customizable RGB lighting', 129.99, 150, 'https://images.unsplash.com/photo-1517336714731-489689fd1ca8', 4.6, 95, NOW(), NOW()),
(1, 'Smartwatch Active 5', 'Fitness tracking smartwatch with heart-rate monitor', 249.99, 120, 'https://images.unsplash.com/photo-1512499617640-c2f999098c01', 4.2, 75, NOW(), NOW()),
(1, 'Bluetooth Speaker Boom', 'Portable waterproof Bluetooth speaker', 89.99, 300, 'https://images.unsplash.com/photo-1544716278-ca5e3f4abd8c', 4.1, 65, NOW(), NOW()),

-- Books
(2, 'Go Programming Essentials', 'Beginner-friendly guide to Golang', 29.99, 100, 'https://images.unsplash.com/photo-1553729784-e91953dec042', 4.8, 50, NOW(), NOW()),
(2, 'Machine Learning Fundamentals', 'Comprehensive ML concepts and practices', 49.99, 70, 'https://images.unsplash.com/photo-1553729784-e91953dec042', 4.6, 40, NOW(), NOW()),
(2, 'The Art of Clean Architecture', 'Best practices for designing scalable systems', 39.99, 60, 'https://images.unsplash.com/photo-1507842217343-583bb7270b66', 4.9, 35, NOW(), NOW()),
(2, 'History of Modern Japan', 'Insightful book about Japan’s modernization', 24.99, 40, 'https://images.unsplash.com/photo-1521587760476-6c12a4b040da', 4.4, 20, NOW(), NOW()),
(2, 'Fantasy Saga: Dragon Realm', 'Epic fantasy adventure novel', 19.99, 120, 'https://images.unsplash.com/photo-1507842217343-583bb7270b66', 4.3, 70, NOW(), NOW()),
(2, 'Self-Discipline Mastery', 'Guide to building strong habits and focus', 21.99, 90, 'https://images.unsplash.com/photo-1553729784-e91953dec042', 4.2, 45, NOW(), NOW()),

-- Clothing
(3, 'Classic White T-Shirt', '100% cotton unisex T-shirt', 14.99, 500, 'https://images.unsplash.com/photo-1521572163474-6864f9cf17ab', 4.0, 30, NOW(), NOW()),
(3, 'Blue Denim Jeans', 'Slim fit stretchable denim jeans', 49.99, 200, 'https://images.unsplash.com/photo-1582418702059-97ebafb35d09', 4.5, 55, NOW(), NOW()),
(3, 'Black Hoodie', 'Comfortable fleece-lined hoodie', 39.99, 250, 'https://images.unsplash.com/photo-1551024601-bec78aea704b', 4.7, 110, NOW(), NOW()),
(3, 'Running Shoes Pro', 'Lightweight running shoes with cushioned sole', 79.99, 150, 'https://images.unsplash.com/photo-1528701800489-20be9c7e6d8f', 4.6, 95, NOW(), NOW()),
(3, 'Summer Dress Floral', 'Light summer dress with floral patterns', 34.99, 180, 'https://images.unsplash.com/photo-1520975918313-60fdbd8d9b8c', 4.4, 60, NOW(), NOW()),
(3, 'Winter Jacket HeavyDuty', 'Insulated winter jacket for cold weather', 129.99, 80, 'https://images.unsplash.com/photo-1520975918313-60fdbd8d9b8c', 4.8, 75, NOW(), NOW()),
(3, 'Baseball Cap Classic', 'Adjustable unisex baseball cap', 12.99, 300, 'https://images.unsplash.com/photo-1520975918313-60fdbd8d9b8c', 4.1, 25, NOW(), NOW());

-- +goose Down
DELETE FROM products;
DELETE FROM categories;

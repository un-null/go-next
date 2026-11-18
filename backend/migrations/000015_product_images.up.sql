-- Every product gets its own photo. The seed reused one URL across different
-- products, and seven of the twenty had gone stale, so the catalogue showed a
-- grey placeholder for them.
--
-- Keyed by id because 000013 has already run against the live database; editing
-- the seed would only have changed a fresh install.
UPDATE products SET image_url = 'https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?w=1200&q=80' WHERE id = 1;
UPDATE products SET image_url = 'https://images.unsplash.com/photo-1572569511254-d8f925fe2cbb?w=1200&q=80' WHERE id = 2;
UPDATE products SET image_url = 'https://images.unsplash.com/photo-1603302576837-37561b2e2302?w=1200&q=80' WHERE id = 3;
UPDATE products SET image_url = 'https://images.unsplash.com/photo-1527443224154-c4a3942d3acf?w=1200&q=80' WHERE id = 4;
UPDATE products SET image_url = 'https://images.unsplash.com/photo-1595044426077-d36d9236d54a?w=1200&q=80' WHERE id = 5;
UPDATE products SET image_url = 'https://images.unsplash.com/photo-1546868871-7041f2a55e12?w=1200&q=80' WHERE id = 6;
UPDATE products SET image_url = 'https://images.unsplash.com/photo-1608043152269-423dbba4e7e1?w=1200&q=80' WHERE id = 7;
UPDATE products SET image_url = 'https://images.unsplash.com/photo-1544716278-ca5e3f4abd8c?w=1200&q=80' WHERE id = 8;
UPDATE products SET image_url = 'https://images.unsplash.com/photo-1555949963-aa79dcee981c?w=1200&q=80' WHERE id = 9;
UPDATE products SET image_url = 'https://images.unsplash.com/photo-1481627834876-b7833e8f5570?w=1200&q=80' WHERE id = 10;
UPDATE products SET image_url = 'https://images.unsplash.com/photo-1526481280693-3bfa7568e0f3?w=1200&q=80' WHERE id = 11;
UPDATE products SET image_url = 'https://images.unsplash.com/photo-1544947950-fa07a98d237f?w=1200&q=80' WHERE id = 12;
UPDATE products SET image_url = 'https://images.unsplash.com/photo-1506126613408-eca07ce68773?w=1200&q=80' WHERE id = 13;
UPDATE products SET image_url = 'https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w=1200&q=80' WHERE id = 14;
UPDATE products SET image_url = 'https://images.unsplash.com/photo-1542272604-787c3835535d?w=1200&q=80' WHERE id = 15;
UPDATE products SET image_url = 'https://images.unsplash.com/photo-1556821840-3a63f95609a7?w=1200&q=80' WHERE id = 16;
UPDATE products SET image_url = 'https://images.unsplash.com/photo-1552346154-21d32810aba3?w=1200&q=80' WHERE id = 17;
UPDATE products SET image_url = 'https://images.unsplash.com/photo-1572804013309-59a88b7e92f1?w=1200&q=80' WHERE id = 18;
UPDATE products SET image_url = 'https://images.unsplash.com/photo-1544923246-77307dd654cb?w=1200&q=80' WHERE id = 19;
UPDATE products SET image_url = 'https://images.unsplash.com/photo-1521369909029-2afed882baee?w=1200&q=80' WHERE id = 20;

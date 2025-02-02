-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS ecom.products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    Name TEXT NOT NULL,
    Description TEXT NOT NULL,
    Image TEXT NOT NULL UNIQUE,
    Price INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL
);

INSERT INTO ecom.products (Name, Description, Image, Price) VALUES
    ('Classic White Sneakers', 'Comfortable everyday sneakers with premium cushioning', '/images/white-sneakers-001.jpg', 7999),
    ('Leather Messenger Bag', 'Handcrafted genuine leather bag with multiple compartments', '/images/leather-bag-002.jpg', 12999),
    ('Wireless Headphones', 'Premium noise-canceling headphones with 30-hour battery life', '/images/headphones-003.jpg', 24999),
    ('Smart Watch Series X', 'Feature-rich smartwatch with health tracking capabilities', '/images/smartwatch-004.jpg', 29999),
    ('Organic Cotton T-Shirt', 'Sustainable, soft cotton t-shirt in classic fit', '/images/tshirt-005.jpg', 2499);

-- Create the function to update the updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = CURRENT_TIMESTAMP;
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create the trigger to call the function on UPDATE
CREATE TRIGGER set_updated_at
BEFORE UPDATE ON ecom.products
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS ecom.products;
-- +goose StatementEnd

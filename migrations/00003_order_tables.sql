-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
SELECT 'up SQL query';
CREATE TABLE ecom.orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    total FLOAT NOT NULL,
    status VARCHAR(50) NOT NULL,
    address TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL
);

CREATE TABLE ecom.order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL,
    product_id UUID NOT NULL,
    quantity INT NOT NULL,
    price FLOAT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    CONSTRAINT fk_order
        FOREIGN KEY(order_id) 
        REFERENCES ecom.orders(id),
    CONSTRAINT fk_product
        FOREIGN KEY(product_id)
        REFERENCES ecom.products(id)
);

-- Create the function to update the updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = CURRENT_TIMESTAMP;
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create the trigger to call the function on UPDATE for orders
CREATE TRIGGER set_updated_at_orders
BEFORE UPDATE ON ecom.orders
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Create the trigger to call the function on UPDATE for order_items
CREATE TRIGGER set_updated_at_order_items
BEFORE UPDATE ON ecom.order_items
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS ecom.order_items;
DROP TABLE IF EXISTS ecom.orders;
-- +goose StatementEnd

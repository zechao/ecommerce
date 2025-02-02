-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

-- Insert example data into products
INSERT INTO ecom.products (id, Name, Description, Image, Price,Quantity)
VALUES 
    ('55555555-5555-5555-5555-555555555555', 'Product 1', 'Description for product 1', 'image1.png', 100, 10),
    ('66666666-6666-6666-6666-666666666666', 'Product 2', 'Description for product 2', 'image2.png', 200, 10),
    ('77777777-7777-7777-7777-777777777777', 'Product 3', 'Description for product 3', 'image3.png', 300, 10),
    ('88888888-8888-8888-8888-888888888888', 'Product 4', 'Description for product 4', 'image4.png', 400, 10),
    ('99999999-9999-9999-9999-999999999999', 'Product 5', 'Description for product 5', 'image5.png', 500, 10),
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Product 6', 'Description for product 6', 'image6.png', 600, 10),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'Product 7', 'Description for product 7', 'image7.png', 700, 10),
    ('cccccccc-cccc-cccc-cccc-cccccccccccc', 'Product 8', 'Description for product 8', 'image8.png', 800, 10),
    ('dddddddd-dddd-dddd-dddd-dddddddddddd', 'Product 9', 'Description for product 9', 'image9.png', 900, 10),
    ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'Product 10', 'Description for product 10', 'image10.png', 1000, 10),
    ('157d8994-121c-4951-a1f1-c6dee8c5835b', 'Classic White Sneakers', 'Comfortable everyday sneakers with premium cushioning', '/images/white-sneakers-001.jpg', 7999, 10),
    ('463ef401-ee47-423d-868d-693fc48c67a4', 'Leather Messenger Bag', 'Handcrafted genuine leather bag with multiple compartments', '/images/leather-bag-002.jpg', 12999, 10),
    ('1898c328-3548-4732-8242-91872703f2b5', 'Wireless Headphones', 'Premium noise-canceling headphones with 30-hour battery life', '/images/headphones-003.jpg', 24999, 10),
    ('ac8c728b-3efc-44eb-95c2-b918e15caa76', 'Smart Watch Series X', 'Feature-rich smartwatch with health tracking capabilities', '/images/smartwatch-004.jpg', 29999, 10),
    ('d21163b0-457b-41b0-ba45-56af5dfec222', 'Organic Cotton T-Shirt', 'Sustainable, soft cotton t-shirt in classic fit', '/images/tshirt-005.jpg', 2499, 10);



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';



-- Delete example data from products
DELETE FROM ecom.products WHERE id IN ('55555555-5555-5555-5555-555555555555', '66666666-6666-6666-6666-666666666666', '77777777-7777-7777-7777-777777777777', '88888888-8888-8888-8888-888888888888', '99999999-9999-9999-9999-999999999999', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 'dddddddd-dddd-dddd-dddd-dddddddddddd', 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', '157d8994-121c-4951-a1f1-c6dee8c5835b', '463ef401-ee47-423d-868d-693fc48c67a4', '1898c328-3548-4732-8242-91872703f2b5', 'ac8c728b-3efc-44eb-95c2-b918e15caa76', 'd21163b0-457b-41b0-ba45-56af5dfec222');


-- +goose StatementEnd
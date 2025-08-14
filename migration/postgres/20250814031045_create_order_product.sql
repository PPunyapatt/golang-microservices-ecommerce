-- +goose Up
-- +goose StatementBegin
CREATE TABLE order_products (
    product_id INT PRIMARY KEY,
    product_name VARCHAR(255),
    price DECIMAL(10,2),
    store_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_products;
-- +goose StatementEnd

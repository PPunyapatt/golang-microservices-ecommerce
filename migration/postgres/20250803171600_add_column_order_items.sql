-- +goose Up
-- +goose StatementBegin
ALTER TABLE order_items
ADD COLUMN store_id INT NOT NULL,
ADD COLUMN product_name VARCHAR(100) NOT NULL,
ADD COLUMN unit_price DECIMAL(10,5) NOT NULL;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE order_items RENAME COLUMN qty TO quantity;
ALTER TABLE order_items RENAME COLUMN price TO total_price;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE order_items
ALTER COLUMN total_price TYPE DECIMAL(10,6);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE order_items
ALTER COLUMN price TYPE DECIMAL(10,2);

ALTER TABLE order_items RENAME COLUMN quantity TO qty;
ALTER TABLE order_items RENAME COLUMN total_price TO price;

ALTER TABLE order_items
DROP COLUMN store_id,
DROP COLUMN product_name,
DROP COLUMN unit_price;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
ADD COLUMN payment_status VARCHAR(40),
ADD COLUMN shipping_address_id INT;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE orders
RENAME COLUMN amount TO total_amount;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders
RENAME COLUMN total_amount TO amount;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE orders
DROP COLUMN payment_status,
DROP COLUMN shipping_address_id;
-- +goose StatementEnd

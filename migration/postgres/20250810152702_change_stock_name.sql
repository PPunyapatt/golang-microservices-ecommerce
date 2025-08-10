-- +goose Up
-- +goose StatementBegin
ALTER TABLE products
RENAME COLUMN stock to available_stock;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE products
ADD COLUMN reserved_stock INT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE products
RENAME COLUMN available_stock to stock;

ALTER TABLE products
DROP COLUMN reserved_stock;
-- +goose StatementEnd

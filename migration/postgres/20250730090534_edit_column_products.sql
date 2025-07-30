-- +goose Up
-- +goose StatementBegin
ALTER TABLE products RENAME COLUMN catagory_id TO category_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE products RENAME COLUMN category_id TO catagory_id;
-- +goose StatementEnd

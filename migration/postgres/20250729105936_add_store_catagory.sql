-- +goose Up
-- +goose StatementBegin
ALTER TABLE categories
ADD COLUMN store_id INT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE categories DROP COLUMN store_id;
-- +goose StatementEnd

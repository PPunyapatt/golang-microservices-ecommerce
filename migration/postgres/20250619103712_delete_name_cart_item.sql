-- +goose Up
-- +goose StatementBegin
ALTER TABLE cart_item DROP COLUMN image_url;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
AltER TABLE cart_item ADD COLUMN image_url VARCHAR;
-- +goose StatementEnd

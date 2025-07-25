-- +goose Up
-- +goose StatementBegin
ALTER TABLE stores DROP CONSTRAINT stores_id_owner_unique;
ALTER TABLE stores ADD CONSTRAINT stores_id_owner_unique UNIQUE (owner);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd

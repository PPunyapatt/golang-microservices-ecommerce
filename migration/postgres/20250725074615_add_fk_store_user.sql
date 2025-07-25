-- +goose Up
-- +goose StatementBegin
ALTER TABLE stores ADD CONSTRAINT fk_owner_user FOREIGN KEY (owner) REFERENCES users(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE stores
DROP CONSTRAINT IF EXISTS fk_owner_user;
-- +goose StatementEnd

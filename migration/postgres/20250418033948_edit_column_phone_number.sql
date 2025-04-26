-- +goose Up
-- +goose StatementBegin
ALTER TABLE users RENAME COLUMN phone_numer TO phone_number;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users RENAME COLUMN phone_number TO phone_numer;
-- +goose StatementEnd

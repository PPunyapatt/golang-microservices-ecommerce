-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD constraint users_email_unique UNIQUE (email);
ALTER TABLE users ADD constraint users_first_last_name_unique UNIQUE (first_name, last_name);
ALTER TABLE users ADD constraint users_phone_number_unique UNIQUE (phone_number);
ALTER TABLE users ADD constraint users_id_unique UNIQUE (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP CONSTRAINT users_email_unique;
ALTER TABLE users DROP CONSTRAINT users_first_last_name_unique;
ALTER TABLE users DROP CONSTRAINT users_phone_number_unique;
ALTER TABLE users DROP CONSTRAINT users_id_unique;
-- +goose StatementEnd

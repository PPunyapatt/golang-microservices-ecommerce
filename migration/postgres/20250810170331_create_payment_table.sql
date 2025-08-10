-- +goose Up
-- +goose StatementBegin
CREATE TABLE payments (
    payment_id VARCHAR(40) PRIMARY KEY,
    order_id INT NOT NULL,
    amount DECIMAL(10,2),
    currency VARCHAR(3),
    payment_method VARCHAR,
    status VARCHAR,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    failure_reason TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payments;
-- +goose StatementEnd

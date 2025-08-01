-- +goose Up
-- +goose StatementBegin
ALTER TABLE cart_items
ADD CONSTRAINT cart_items_cartid_productid_unique UNIQUE (cart_id, product_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE cart_items
DROP CONSTRAINT cart_items_cartid_productid_unique;
-- +goose StatementEnd

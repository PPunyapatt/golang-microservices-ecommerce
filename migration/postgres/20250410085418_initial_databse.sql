-- +goose Up
-- +goose StatementBegin
-- ------------------------------- User -------------------------------
CREATE TABLE users (
  id UUID PRIMARY KEY,
  first_name VARCHAR(50) NOT NULL,
  last_name VARCHAR(50) NOT NULL,
  email VARCHAR(50) NOT NULL,
  password_hash VARCHAR(100) NOT NULL,
  phone_number VARCHAR NOT NULL,
  verified BOOLEAN,
  image_url VARCHAR,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP
);

CREATE TABLE roles (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL
);

CREATE TABLE user_role (
  user_id UUID NOT NULL,
  role_id INT NOT NULL,
  PRIMARY KEY (user_id, role_id),
  CONSTRAINT fk_user_role FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_role_role FOREIGN KEY (role_id) REFERENCES roles(id)
);

CREATE TABLE address (
  id SERIAL PRIMARY KEY,
  user_id UUID NOT NULL,
  address VARCHAR(100) NOT NULL,
  city VARCHAR(20) NOT NULL,
  country VARCHAR(20) NOT NULL,
  post_code VARCHAR(5) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP,
  CONSTRAINT fk_user_address FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE bank_account (
  id SERIAL PRIMARY KEY,
  user_id UUID NOT NULL,
  account_number VARCHAR NOT NULL,
  bank_name VARCHAR NOT NULL,
  payment_type VARCHAR NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP,
  CONSTRAINT fk_user_bank FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE stores (
  id SERIAL PRIMARY KEY,
  name VARCHAR(1000) NOT NULL,
  owner UUID NOT NULL,
  CONSTRAINT fk_owner_user FOREIGN KEY (owner) REFERENCES users(id),
  CONSTRAINT stores_id_owner_unique UNIQUE (owner)
);

-- ------------------------------- Inventory -------------------------------
CREATE TABLE catagories (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  store_id INT NOT NULL
);

CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  name VARCHAR,
  description TEXT,
  store_id INT NOT NULL,
  category_id INT,
  image_url VARCHAR,
  price DECIMAL(10,2),
  stock INT DEFAULT 0,
  add_by UUID NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  CONSTRAINT fk_catagory_product FOREIGN KEY (category_id) REFERENCES catagories(id) ON DELETE SET NULL
);

-- ------------------------------- Order -------------------------------
CREATE TABLE orders (
  id SERIAL PRIMARY KEY,
  user_id UUID NOT NULL,
  status VARCHAR(20),
  amount INT,
  payment_id INT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

CREATE TABLE order_items (
  id SERIAL PRIMARY KEY,
  order_id INT NOT NULL,
  product_id INT NOT NULL,
  qty INT NOT NULL,
  price DECIMAL(10, 2) NOT NULL
);

CREATE TABLE carts (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE INDEX idx_cart_user_product ON carts (user_id);

CREATE TABLE cart_items (
    id SERIAL PRIMARY KEY,
    cart_id INT NOT NULL,
    product_id INT NOT NULL,
    product_name VARCHAR(200) NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    price DECIMAL(10, 5) NOT NULL,
    store_id INT NOT NULL,
    image_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (cart_id) REFERENCES carts(id) ON DELETE CASCADE,
    CONSTRAINT cart_items_cartid_productid_unique UNIQUE (cart_id, product_id)
);

CREATE INDEX idx_cart_items_cart ON cart_items (cart_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cart;
DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS catagories;
DROP TABLE IF EXISTS bank_account;
DROP TABLE IF EXISTS address;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd

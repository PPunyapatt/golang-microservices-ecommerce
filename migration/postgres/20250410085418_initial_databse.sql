
-- +goose Up
CREATE TABLE users (
  id CHAR(36) PRIMARY KEY,
  first_name VARCHAR NOT NULL,
  last_name VARCHAR NOT NULL,
  email VARCHAR NOT NULL,
  phone_numer VARCHAR NOT NULL,
  verified BOOLEAN,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

CREATE TABLE role (
  id SERIAL PRIMARY KEY,
  name VARCHAR NOT NULL
);

CREATE TABLE user_role (
  user_id CHAR(36) NOT NULL,
  role_id INT NOT NULL,
  CONSTRAINT fk_user_role FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_role_role FOREIGN KEY (role_id) REFERENCES role(id)
);


CREATE TABLE address (
  id SERIAL PRIMARY KEY,
  user_id CHAR(36) NOT NULL,
  address VARCHAR NOT NULL,
  city VARCHAR NOT NULL,
  country VARCHAR NOT NULL,
  post_code VARCHAR NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  CONSTRAINT fk_user_address FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE bank_account (
  id SERIAL PRIMARY KEY,
  user_id CHAR(36) NOT NULL,
  bank_account VARCHAR NOT NULL,
  payment_type VARCHAR NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  CONSTRAINT fk_user_bank FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE catagory (
  id SERIAL PRIMARY KEY,
  name VARCHAR,
  image_url VARCHAR,
  display_order INT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

CREATE TABLE product (
  id SERIAL PRIMARY KEY,
  name VARCHAR,
  description TEXT,
  catagory_id INT,
  image_url VARCHAR,
  price FLOAT,
  stock INT,
  user_id CHAR(36) NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  CONSTRAINT fk_user_product FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_catagory_product FOREIGN KEY (catagory_id) REFERENCES catagory(id) ON DELETE SET NULL
);

CREATE TABLE payment (
  id SERIAL PRIMARY KEY,
  user_id CHAR(36) NOT NULL,
  capture_method VARCHAR,
  amount FLOAT,
  transaction_id INT,
  customer_id INT,
  payment_id INT,
  status VARCHAR,
  response VARCHAR,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  CONSTRAINT fk_user_payment FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE orders (
  id SERIAL PRIMARY KEY,
  user_id CHAR(36) NOT NULL,
  status VARCHAR,
  amount INT,
  transaction_id INT,
  order_ref_number INT,
  payment_id INT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  CONSTRAINT fk_user_order FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_payment_order FOREIGN KEY (payment_id) REFERENCES payment(id) ON DELETE SET NULL
);

CREATE TABLE order_item (
  id SERIAL PRIMARY KEY,
  order_id INT NOT NULL,
  product_id INT,
  name VARCHAR,
  image_url VARCHAR,
  seller_id VARCHAR,
  price FLOAT,
  qty INT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  CONSTRAINT fk_order_item_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_product_order_item FOREIGN KEY (product_id) REFERENCES product(id) ON DELETE SET NULL
);

CREATE TABLE cart (
  id SERIAL PRIMARY KEY,
  user_id CHAR(36) NOT NULL,
  product_id INT NOT NULL,
  name VARCHAR NOT NULL,
  image_url VARCHAR,
  price FLOAT NOT NULL,
  qty INT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  CONSTRAINT fk_user_cart FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_product_cart FOREIGN KEY (product_id) REFERENCES product(id) ON DELETE CASCADE
);


CREATE TABLE shipping (
  id SERIAL PRIMARY KEY,
  order_id INT NOT NULL,
  buyer_id VARCHAR NOT NULL,
  seller_id VARCHAR NOT NULL,
  status VARCHAR,
  shipping_address TEXT,
  CONSTRAINT fk_shipping_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_shipping_buyer FOREIGN KEY (buyer_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_shipping_seller FOREIGN KEY (seller_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS shipping;
DROP TABLE IF EXISTS cart;
DROP TABLE IF EXISTS order_item;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS payment;
DROP TABLE IF EXISTS product;
DROP TABLE IF EXISTS catagory;
DROP TABLE IF EXISTS bank_account;
DROP TABLE IF EXISTS address;
DROP TABLE IF EXISTS users;


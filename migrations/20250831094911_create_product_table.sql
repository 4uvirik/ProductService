-- +goose Up
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    category_id INT REFERENCES category(id) ON DELETE SET NULL
);

-- +goose Down
DROP TABLE IF EXISTS products;


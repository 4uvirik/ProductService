-- +goose Up
CREATE TABLE IF NOT EXISTS category (
    id SERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS category;

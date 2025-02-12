-- +goose Up
CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    price INT NOT NULL
);

-- +goose Down
-- DROP TABLE IF EXISTS items;

-- +goose Up
CREATE TABLE IF NOT EXISTS inventories (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_type VARCHAR(255) NOT NULL,
    quantity INT NOT NULL CHECK (quantity >= 0)
);

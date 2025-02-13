-- +goose Up
CREATE TABLE IF NOT EXISTS user_inventory (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    item_id INT REFERENCES items(id),
    quantity INT,
    CONSTRAINT unique_item_per_user UNIQUE (user_id, item_id)
);

-- -- +goose Down
-- DROP TABLE IF EXISTS user_inventory;
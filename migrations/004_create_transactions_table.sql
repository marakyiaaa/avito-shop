-- +goose Up
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    from_user_id INT REFERENCES users(id),
    to_user_id INT REFERENCES users(id),
    amount INT NOT NULL CHECK (amount > 0),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- -- +goose Down
-- DROP TABLE IF EXISTS transactions;
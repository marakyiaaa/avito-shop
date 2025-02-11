-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
                                     id       SERIAL PRIMARY KEY,
                                     username VARCHAR(255) UNIQUE NOT NULL,
                                     password VARCHAR(255)        NOT NULL,
                                     coins    INT DEFAULT 1000
);

-- +migrate Down
DROP TABLE IF EXISTS users;
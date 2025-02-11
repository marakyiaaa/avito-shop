-- +migrate Up
CREATE TABLE IF NOT EXISTS items (
                                     id    SERIAL PRIMARY KEY,
                                     name  VARCHAR(255) UNIQUE NOT NULL,
                                     price INT                 NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS items;
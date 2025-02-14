-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id       SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255)        NOT NULL,
    balance    INT DEFAULT 1000
);

-- INSERT INTO goose_db_version (version_id, is_applied)
-- SELECT 1, TRUE
-- WHERE NOT EXISTS (SELECT 1 FROM goose_db_version WHERE version_id = 1);
-- !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
-- !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
-- !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
-- !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!


-- -- +goose Down
-- DROP TABLE IF EXISTS users;
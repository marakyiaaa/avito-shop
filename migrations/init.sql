CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(255) UNIQUE NOT NULL,
                       password VARCHAR(255) NOT NULL,
                       coins INT DEFAULT 1000
);

CREATE TABLE items (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(255) UNIQUE NOT NULL,
                       price INT NOT NULL
);

CREATE TABLE transactions (
                              id SERIAL PRIMARY KEY,
                              from_user_id INT REFERENCES users(id),
                              to_user_id INT REFERENCES users(id),
                              amount INT NOT NULL,
                              timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- goose
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username varchar(32) NOT NULL UNIQUE,
    password_hash TEXT
)
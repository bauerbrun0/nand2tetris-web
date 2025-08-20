CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    email_verified BOOLEAN,
    password_hash CHAR(97),
    created TIMESTAMPTZ NOT NULL,
    CONSTRAINT users_unique_constraint_username UNIQUE (username),
    CONSTRAINT users_unique_constraint_email UNIQUE (email)
);

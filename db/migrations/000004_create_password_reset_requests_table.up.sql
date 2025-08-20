CREATE TABLE IF NOT EXISTS password_reset_requests (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    email VARCHAR(255) NOT NULL,
    code CHAR(12) NOT NULL,
    verify_email_after BOOLEAN,
    expiry TIMESTAMPTZ NOT NULL,
    CONSTRAINT password_reset_requests_unique_constraint_code UNIQUE (code),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

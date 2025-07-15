CREATE TABLE IF NOT EXISTS email_verification_requests (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    email VARCHAR(255) NOT NULL,
    code CHAR(8) NOT NULL,
    expiry TIMESTAMPTZ NOT NULL,
    CONSTRAINT email_verification_requests_unique_constraint_user_id UNIQUE (user_id),
    CONSTRAINT email_verification_requests_unique_constraint_code UNIQUE (code),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

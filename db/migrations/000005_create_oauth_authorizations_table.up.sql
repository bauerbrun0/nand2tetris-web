CREATE TYPE provider AS ENUM ('GitHub', 'Google');

CREATE TABLE IF NOT EXISTS oauth_authorizations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    provider provider NOT NULL,
    user_provider_id TEXT NOT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT oauth_authorizations_unique_constraint UNIQUE (user_id, provider, user_provider_id)
);

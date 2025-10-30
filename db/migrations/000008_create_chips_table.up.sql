CREATE TABLE IF NOT EXISTS chips (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    hdl TEXT,
    created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chips_unique_constraint_project_id_name UNIQUE (project_id, name),
    CONSTRAINT fk_project_id FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE
);

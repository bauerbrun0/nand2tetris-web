-- create nand2tetris database
CREATE DATABASE nand2tetris_web;

-- connect to db
\c nand2tetris_web

-- create migration user
CREATE USER nand2tetris_web_migration WITH PASSWORD 'password';

-- grant privileges to migration user
GRANT ALL PRIVILEGES ON DATABASE nand2tetris_web TO nand2tetris_web_migration;
ALTER DEFAULT PRIVILEGES GRANT ALL ON TABLES TO nand2tetris_web_migration;
ALTER DEFAULT PRIVILEGES GRANT ALL ON SEQUENCES TO nand2tetris_web_migration;
-- set as owner of the schema
ALTER SCHEMA public OWNER TO nand2tetris_web_migration;

-- create application user with more limited permissions
CREATE USER nand2tetris_web WITH PASSWORD 'password';

-- setting default privileges with postgres user would have no effect
-- since the tables are created with the migration user
SET ROLE nand2tetris_web_migration;

-- grant basic DML operations on all existing tables
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO nand2tetris_web;
-- grant usage and select on seuences
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO nand2tetris_web;
-- grant the same for future tables/sequences
ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO nand2tetris_web;
ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT USAGE, SELECT ON SEQUENCES TO nand2tetris_web;

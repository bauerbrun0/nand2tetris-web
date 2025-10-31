DROP TRIGGER IF EXISTS trg_projects_generate_slug ON projects;
DROP FUNCTION IF EXISTS projects_generate_slug();
DROP FUNCTION IF EXISTS generate_unique_slug(INT, TEXT);
DROP FUNCTION IF EXISTS slugify(TEXT);
DROP EXTENSION IF EXISTS unaccent CASCADE;

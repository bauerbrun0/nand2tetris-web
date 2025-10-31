CREATE EXTENSION IF NOT EXISTS unaccent;

CREATE OR REPLACE FUNCTION slugify(input TEXT)
RETURNS TEXT AS $$
DECLARE
    slug TEXT;
BEGIN
    slug := lower(trim(regexp_replace(unaccent(input), '[^a-zA-Z0-9]+', '-', 'g')));
    slug := trim(both '-' FROM slug);
    RETURN slug;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

CREATE OR REPLACE FUNCTION generate_unique_slug(p_user_id INT, p_title TEXT)
RETURNS TEXT AS $$
DECLARE
    base_slug TEXT := slugify(p_title);
    new_slug TEXT := base_slug;
    counter INT := 1;
BEGIN
    -- loop until we find a slug that doesn't exist for this user
    WHILE EXISTS (
        SELECT 1 FROM projects
        WHERE projects.user_id = p_user_id
          AND projects.slug = new_slug
    ) LOOP
        counter := counter + 1;
        new_slug := base_slug || '-' || counter;
    END LOOP;

    RETURN new_slug;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION projects_generate_slug()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        -- always generate slug on insert
        NEW.slug := generate_unique_slug(NEW.user_id, NEW.title);

    ELSIF TG_OP = 'UPDATE' THEN
        -- only regenerate slug if title changed
        IF NEW.title <> OLD.title THEN
            NEW.slug := generate_unique_slug(NEW.user_id, NEW.title);
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_projects_generate_slug
BEFORE INSERT OR UPDATE ON projects
FOR EACH ROW
EXECUTE FUNCTION projects_generate_slug();

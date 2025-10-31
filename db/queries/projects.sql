-- name: CreateProject :one
INSERT INTO projects (
    user_id, title, description
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetProject :one
SELECT
    id, user_id, title, slug, description, created, updated
FROM projects
WHERE user_id = $1 AND id = $2;

-- name: GetProjectBySlug :one
SELECT
    id, user_id, title, slug, description, created, updated
FROM projects
WHERE user_id = $1 AND slug = $2;

-- name: UpdateProject :one
UPDATE projects SET
    title = $3, description = $4, updated = NOW()
WHERE user_id = $1 AND id = $2
RETURNING *;

-- name: DeleteProject :one
DELETE FROM projects
WHERE user_id = $1 AND id = $2
RETURNING *;

-- name: GetPaginatedProjects :many
SELECT
    id, user_id, title, slug, description, created, updated
FROM projects
WHERE user_id = $1
ORDER BY updated DESC
LIMIT sqlc.arg(pagelimit)::integer
OFFSET sqlc.arg(pageoffset)::integer;

-- name: GetProjectsCount :one
SELECT COUNT(*) AS count FROM projects
WHERE user_id = $1;

-- name: IsProjectOwnedByUser :one
SELECT EXISTS (
    SELECT 1 FROM projects
    WHERE id = $1 AND user_id = $2
) AS exists;

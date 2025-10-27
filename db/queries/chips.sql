-- name: CreateChip :one
INSERT INTO chips (
    project_id, name
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetChip :one
SELECT
    id, project_id, name, hdl, created, updated
FROM chips
WHERE id = $1 AND project_id = $2;

-- name: GetChipsByProject :many
SELECT
    id, project_id, name, hdl, created, updated
FROM chips
WHERE project_id = $1
ORDER BY name ASC;

-- name: UpdateChip :one
UPDATE chips SET
    name = $2, hdl = $3, updated = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteChip :one
DELETE FROM chips
WHERE id = $1
RETURNING *;

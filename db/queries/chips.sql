-- name: CreateChip :one
INSERT INTO chips (
    project_id, name
) VALUES (
    $1, $2
) RETURNING *;

-- name: IsChipOwnedByUser :one
SELECT EXISTS (
    SELECT 1
    FROM chips c
    JOIN projects p ON c.project_id = p.id
    WHERE c.id = $1 AND p.user_id = $2
) AS exists;

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

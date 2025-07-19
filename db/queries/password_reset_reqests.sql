-- name: CreatePasswordResetRequest :one
INSERT INTO password_reset_requests (
    user_id, email, code, expiry, verify_email_after
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetPasswordResetRequestByCode :one
SELECT * FROM password_reset_requests WHERE
    code = $1;

-- name: InvalidatePasswordResetRequest :exec
UPDATE password_reset_requests
SET expiry = sqlc.arg(now)::timestamptz
WHERE
    id = $1;

-- name: InvalidatePasswordResetRequestsOfUser :exec
UPDATE password_reset_requests
SET expiry = sqlc.arg(now)::timestamptz
WHERE
    user_id = $1 AND
    expiry > sqlc.arg(now)::timestamptz;

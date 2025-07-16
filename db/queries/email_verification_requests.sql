-- name: CreateEmailVerificationRequest :one
INSERT INTO email_verification_requests (
    user_id, email, code, expiry
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetEmailVerificationRequestByCode :one
SELECT * FROM email_verification_requests WHERE
    code = $1;

-- name: InvalidateEmailVerificationRequest :exec
UPDATE email_verification_requests
SET expiry = sqlc.arg(now)::timestamptz
WHERE
    id = $1;

-- name: InvalidateEmailVerificationRequestsOfUser :exec
UPDATE email_verification_requests
SET expiry = sqlc.arg(now)::timestamptz
WHERE
    user_id = $1 AND
    expiry > sqlc.arg(now)::timestamptz;

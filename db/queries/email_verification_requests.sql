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

-- name: DeleteEmailVerificationRequest :exec
DELETE FROM email_verification_requests WHERE
    id = $1;

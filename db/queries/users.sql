-- name: CreateNewUser :one
INSERT INTO users (
    username, email, email_verified, password_hash, created
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: VerifyUserEmail :exec
UPDATE users SET email_verified = true WHERE
    id = $1;

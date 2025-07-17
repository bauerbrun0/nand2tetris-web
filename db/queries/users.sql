-- name: CreateNewUser :one
INSERT INTO users (
    username, email, email_verified, password_hash, created
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByUsernameOrEmail :one
SELECT * FROM users
WHERE
    LOWER(email) = LOWER(sqlc.arg(identifier)::text) OR
    LOWER(username) = LOWER(sqlc.arg(identifier)::text);

-- name: VerifyUserEmail :exec
UPDATE users SET email_verified = true WHERE
    id = $1;

-- name: GetUserInfo :one
SELECT id, username, email, created FROM users
WHERE
    id = $1;

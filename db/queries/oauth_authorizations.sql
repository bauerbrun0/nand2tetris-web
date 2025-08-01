-- name: CreateOAuthAuthorization :one
INSERT INTO oauth_authorizations (
    user_id, provider, user_provider_id
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: FindOAuthAuthorization :one
SELECT * FROM oauth_authorizations WHERE
    user_provider_id = $1 AND provider = $2;


-- name: DeleteOAuthAuthorization :exec
DELETE FROM oauth_authorizations WHERE
    user_id = $1 AND provider = $2;

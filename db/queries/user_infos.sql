-- name: GetUserInfo :one
SELECT * FROM user_infos
WHERE
    id = $1;

-- name: GetUserInfoByEmailOrUsername :one
SELECT * FROM user_infos
WHERE
    username = $1 OR email = $2;

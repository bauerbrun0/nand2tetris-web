CREATE
OR REPLACE VIEW user_infos AS
SELECT
    u.id,
    u.username,
    u.email,
    u.email_verified,
    u.created,
    u.password_hash IS NOT NULL AS is_password_set,
    COALESCE(
        ARRAY_AGG (a.provider) FILTER (
            WHERE
                a.provider IS NOT NULL
        ),
        ARRAY[]::provider[]
    )::text[] AS linked_accounts
FROM
    users u
    LEFT JOIN oauth_authorizations a ON a.user_id = u.id
GROUP BY
    u.id;

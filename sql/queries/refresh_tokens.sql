-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3
)
RETURNING *;

-- name: SetTokenRevokedAt :one
UPDATE refresh_tokens 
SET revoked_at = $2, updated_at = $2
WHERE token = $1
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- name: GetUserFromRefreshToken :one
SELECT * FROM users
WHERE id = (
    SELECT user_id From refresh_tokens
    WHERE token = $1
);
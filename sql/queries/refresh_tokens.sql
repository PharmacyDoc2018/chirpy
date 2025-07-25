-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES ($1, $2, $2, $3, $4)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
SET expires_at = $2
AND updated_at = $2
WHERE token = $1
RETURNING *;
-- name: UpdatePassword :exec
UPDATE users
SET hashed_password = $2
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: UpdateEmailAndPasswordByID :one
UPDATE users
SET email = $2, hashed_password = $3, updated_at = $4
WHERE id = $1
RETURNING *;
-- name: UpdatePassword :exec
UPDATE users
SET hashed_password = $2
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;
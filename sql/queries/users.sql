-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES($1, $2, $3, $4) RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE name=$1 LIMIT 1;

-- name: DeleteAllUsers :exec
DELETE FROM users WHERE id IS NOT NULL;

-- name: GetAllUsers :many
SELECT * FROM users;


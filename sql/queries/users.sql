-- name: CreateUser :one
INSERT INTO users (id, email) 
VALUES ($1, $2)
RETURNING *;

-- name: DeleteAllUser :exec
DELETE FROM users;
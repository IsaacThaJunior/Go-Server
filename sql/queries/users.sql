-- name: CreateUser :one
INSERT INTO users (id, email, hashed_password) 
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: DeleteAllUser :exec
DELETE FROM users;
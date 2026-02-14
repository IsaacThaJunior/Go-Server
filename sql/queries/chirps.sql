-- name: CreateChirp :one
INSERT INTO chirps (id, body, user_id) 
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM chirps;

-- name: GetChirp :one
SELECT * FROM chirps
WHERE id = $1;

-- name: DeleteAllChirps :exec
DELETE FROM chirps;
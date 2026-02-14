-- name: CreateChirp :one
INSERT INTO chirps (id, body, user_id) 
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteAllChirps :exec
DELETE FROM chirps;
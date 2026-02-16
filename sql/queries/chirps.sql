-- name: CreateChirp :one
INSERT INTO chirps (id, body, user_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAllChirps :many
SELECT *
FROM chirps
ORDER BY created_at ASC;

-- name: GetChirpsByAuthor :many
SELECT * FROM chirps
WHERE user_id = $1
ORDER BY created_at ASC;

-- name: GetChirp :one
SELECT *
FROM chirps
WHERE id = $1;

-- name: DeleteOneChirp :exec
DELETE FROM chirps
WHERE id = $1;

-- name: DeleteAllChirps :exec
DELETE FROM chirps;
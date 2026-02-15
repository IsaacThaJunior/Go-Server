-- name: CreateUser :one
INSERT INTO users (id, email, hashed_password)
VALUES ($1, $2, $3)
RETURNING *;
-- name: GetUser :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET email = $2,
  hashed_password = $3,
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpgradeChirpyStat :one
UPDATE users
SET is_chirpy_red = $2,
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteAllUser :exec
DELETE FROM users;
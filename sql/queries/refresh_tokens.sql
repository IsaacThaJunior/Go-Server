-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, user_id, expires_at) 
VALUES ($1, $2, $3)
RETURNING *;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens 
SET revoked_at = NOW(), updated_at = NOW()
where token = $1;

-- name: GetUserFromRefreshToken :one
SELECT * FROM
refresh_tokens 
WHERE token = $1
AND revoked_at IS NULL
AND expires_at > NOW();
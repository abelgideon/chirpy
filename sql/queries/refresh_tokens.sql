-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3
)
RETURNING *;

-- name: GetUserByRefreshToken :one
SELECT users.*
FROM refresh_tokens
INNER JOIN users ON refresh_tokens.user_id = users.id
WHERE refresh_tokens.token = $1 AND refresh_tokens.revoked_at IS NULL AND NOW() < refresh_tokens.expires_at;

-- name: SetRevokedAtAndUpdatedAt :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;

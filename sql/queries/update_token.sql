-- name: UpdateToken :exec
UPDATE refresh_tokens
SET revoked_at = $2, updated_at = $2
WHERE refresh_tokens.token = $1;
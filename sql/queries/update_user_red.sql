-- name: UpdateUserRed :one
UPDATE users
SET is_chirpy_red = $2, updated_at = $3
WHERE id = $1
RETURNING *;
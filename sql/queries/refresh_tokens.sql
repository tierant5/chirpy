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

-- name: SetRefreshTokenRevoked :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- -- name: GetAllChirps :many
-- SELECT * FROM chirps
-- ORDER BY created_at ASC;
--
-- -- name: DeleteAllChirps :exec
-- DELETE FROM chirps;
--
-- -- name: GetChirpByID :one
-- SELECT * FROM chirps
-- WHERE id = $1
-- LIMIT 1;

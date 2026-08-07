-- name: CreateOAuthState :one
INSERT INTO oauth_states (token, expires_at, user_id, service, channel_id, channel_handle)
VALUES (
  $1,
  NOW() + INTERVAL '5 minutes',
  $2,
  $3,
  $4,
  $5
)
RETURNING *;

-- name: GetOAuthStateFromToken :one
SELECT *
FROM oauth_states
WHERE token = $1;

-- name: DeleteOAuthStateFromToken :exec
DELETE FROM oauth_states
WHERE token = $1;
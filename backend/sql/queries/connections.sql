-- name: CreateConnection :one
INSERT INTO connections (user_id, service, channel_id, channel_handle, access_token, refresh_token, token_expires_at, scopes)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7,
  $8
)
RETURNING *;


-- name: GetConnections :many
SELECT *
FROM connections
WHERE user_id = $1;


-- name: GetConnectionWithID :one
SELECT *
FROM connections
WHERE id = $1;


-- name: ToggleConnection :one
UPDATE connections
SET updated_at = NOW(),
  active = $2
WHERE id = $1
RETURNING *;
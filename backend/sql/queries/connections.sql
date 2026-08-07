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
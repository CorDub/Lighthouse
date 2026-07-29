-- name: CreateMagicLinkToken :one
INSERT INTO magic_link_tokens (token, created_at, updated_at, user_id, name, expires_at)
VALUES (
  $1,
  NOW(),
  NOW(),
  $2,
  $3,
  NOW() + INTERVAL '1 hour'
)
RETURNING *;

-- name: RevokeMagicLinkToken :exec
UPDATE magic_link_tokens
SET updated_at = NOW(),
  revoked_at = NOW()
WHERE token = $1;

-- name: GetMagicLinkTokenByToken :one
SELECT *
FROM magic_link_tokens
WHERE token = $1;

-- name: RevokeMagicLinkTokensFromUser :exec
UPDATE magic_link_tokens
SET updated_at = NOW(),
  revoked_at = NOW()
WHERE user_id = $1
AND revoked_at IS NULL;

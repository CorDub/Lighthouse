--  name: CreateReport :one
INSERT INTO reports (name, tracking_start, tracking_end, created_at, updated_at)
VALUES (
  $1,
  $2,
  $3,
  NOW(),
  NOW()
)
RETURNING *;
-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password, language, role, name)
VALUES (
  gen_random_uuid(),
  NOW(),
  NOW(),
  $1,
  $2,
  $3,
  $4,
  $5
)
RETURNING *;

-- name: CreateUserViaGoogle :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES (
  gen_random_uuid(),
  NOW(),
  NOW(),
  $1
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetUsers :many
SELECT *
FROM users;

-- name: UpdateUserLanguage :exec
UPDATE users
SET updated_at = NOW(),
  language = $2
WHERE id = $1;

-- name: GetCreators :many
SELECT name
FROM users
WHERE role = 'creator'
ORDER BY name ASC;

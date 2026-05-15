-- +goose Up
ALTER TABLE users
  ALTER COLUMN hashed_password DROP NOT NULL,
  ALTER COLUMN hashed_password DROP DEFAULT;

-- +goose Down
ALTER TABLE users
  ALTER COLUMN hashed_password SET NOT NULL,
  ALTER COLUMN hashed_password SET DEFAULT 'unset';
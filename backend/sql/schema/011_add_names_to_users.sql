-- +goose Up
ALTER TABLE users
ADD COLUMN name TEXT;

ALTER TABLE users
ADD CONSTRAINT creator_requires_name
CHECK (role != 'creator' OR name IS NOT NULL);

-- +goose Down
ALTER TABLE users
DROP CONSTRAINT creator_requires_name;

ALTER TABLE users
DROP COLUMN name;
-- +goose Up
ALTER TABLE users
ADD COLUMN role TEXT NOT NULL DEFAULT 'visitor',
ADD CONSTRAINT role_check CHECK (role IN ('visitor', 'agency', 'creator', 'brand'));

-- +goose Down
ALTER TABLE users
DROP COLUMN role;

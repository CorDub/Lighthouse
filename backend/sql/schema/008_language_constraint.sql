-- +goose Up
ALTER TABLE users
ADD CONSTRAINT users_language_check CHECK (language IN ('en', 'es'));

-- +goose Down
ALTER TABLE users
DROP CONSTRAINT users_language_check;
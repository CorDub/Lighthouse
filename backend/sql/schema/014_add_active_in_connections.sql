-- +goose Up
ALTER TABLE connections
ADD COLUMN active BOOLEAN NOT NULL DEFAULT true;

-- +goose Down
ALTER TABLE connections
DROP COLUMN active;
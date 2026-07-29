-- +goose Up
ALTER TABLE magic_link_tokens 
ADD COLUMN name TEXT;

-- +goose Down
ALTER TABLE magic_link_tokens 
DROP COLUMN name;

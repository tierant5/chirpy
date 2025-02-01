-- +goose Up
ALTER TABLE users
ADD COLUMN is_chirpy_red BOOLEAN NOT NULL DEFAULT false;

-- +goose Down
DROP COLUMN IF EXISTS is_chirpy_red;

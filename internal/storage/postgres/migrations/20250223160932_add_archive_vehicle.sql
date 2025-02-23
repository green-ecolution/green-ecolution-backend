-- +goose Up
ALTER TABLE vehicles ADD COLUMN archived_at TIMESTAMP;

-- +goose Down
ALTER TABLE vehicles DROP COLUMN archived_at;

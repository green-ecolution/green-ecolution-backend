-- +goose Up
ALTER TABLE tree_clusters
ADD COLUMN name TEXT;

-- +goose Down
ALTER TABLE tree_clusters
DROP COLUMN name;
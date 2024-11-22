-- +goose Up
-- +goose StatementBegin
CREATE TYPE driver_licence AS ENUM ('B', 'BE', 'C');

ALTER TABLE vehicles
ADD COLUMN model TEXT NOT NULL DEFAULT 'Unknown',
ADD COLUMN driver_licence driver_licence NOT NULL DEFAULT 'B',
ADD COLUMN height FLOAT NOT NULL DEFAULT 0,
ADD COLUMN length FLOAT NOT NULL DEFAULT 0,
ADD COLUMN width FLOAT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE vehicles
DROP COLUMN model,
DROP COLUMN driver_licence,
DROP COLUMN height,
DROP COLUMN length,
DROP COLUMN width;

DROP TYPE IF EXISTS driver_licence;
-- +goose StatementEnd

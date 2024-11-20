-- +goose Up
-- +goose StatementBegin
INSERT INTO sensors (id, status) VALUES ('sensor-1', 'online');
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO sensor_data (sensor_id, data)
VALUES 
  ('sensor-1', '{"temperature": 20.0, "humidity": 0.5, "moisture": 0.75}');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM sensors;
DELETE FROM sensor_data;
-- +goose StatementEnd

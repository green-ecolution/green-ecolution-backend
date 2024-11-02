-- +goose Up
-- +goose StatementBegin
INSERT INTO sensors (id, status) VALUES (1, 'online');
ALTER SEQUENCE sensors_id_seq RESTART WITH 2;
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO sensor_data (sensor_id, data)
VALUES 
  (1, '{"temperature": 20.0, "humidity": 0.5, "moisture": 0.75}');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM sensors;
DELETE FROM sensor_data;
-- +goose StatementEnd

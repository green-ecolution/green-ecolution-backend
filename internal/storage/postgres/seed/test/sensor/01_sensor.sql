-- +goose Up
-- +goose StatementBegin

INSERT INTO sensors (id, status, latitude, longitude, geometry)
VALUES
    ('sensor-1', 'online', 54.82124518093376, 9.485702120628517, ST_SetSRID(ST_MakePoint(54.82124518093376, 9.485702120628517), 4326));

INSERT INTO sensors (id, status, latitude, longitude, geometry)
VALUES
    ('sensor-2', 'offline', 54.78780993841013, 9.444052105200551, ST_SetSRID(ST_MakePoint(54.78780993841013, 9.444052105200551), 4326));

INSERT INTO sensors (id, status, latitude, longitude, geometry)
VALUES
    ('sensor-3', 'unknown', 54.77933725347423, 9.426465409018832, ST_SetSRID(ST_MakePoint(54.77933725347423, 9.426465409018832), 4326));

INSERT INTO sensors (id, status, latitude, longitude, geometry)
VALUES
    ('sensor-4', 'online', 54.82078826498143, 9.489684366114483, ST_SetSRID(ST_MakePoint(54.82078826498143, 9.489684366114483), 4326));

-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO sensor_data (sensor_id, data)
VALUES 
   ('sensor-1', '{"device": "sensor-123", "battery": 34.0, "humidity": 50.0, "temperature": 20.0, "latitude": 54.82124518093376, "longitude": 9.485702120628517, "watermarks":[{"centibar": 38, "resistance": 23, "depth": 30}, {"centibar": 38, "resistance": 23, "depth": 60}, {"centibar": 38, "resistance": 23, "depth": 90}]}');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM sensors;
DELETE FROM sensor_data;
-- +goose StatementEnd

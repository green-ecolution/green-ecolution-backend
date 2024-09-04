-- +goose Up
-- +goose StatementBegin
INSERT INTO images (url) VALUES ('https://avatars.githubusercontent.com/u/165842746?s=96&v=4');
INSERT INTO images (url, filename, mime_type) VALUES ('https://app.dev.green-ecolution.de/api/v1/images/avatar.png', 'avatar.png', 'image/png');

INSERT INTO vehicles (id, number_plate, description, water_capacity) VALUES (1, 'B-1234', 'Test vehicle 1', 100.0);
INSERT INTO vehicles (id, number_plate, description, water_capacity) VALUES (2, 'B-5678', 'Test vehicle 2', 150.0);

INSERT INTO user_vehicles (user_id, vehicle_id) VALUES ('95b69b4c-b38b-4394-9520-496879b67791', 1);
INSERT INTO user_vehicles (user_id, vehicle_id) VALUES ('d2563a8e-a608-4039-8718-25fc3c1d8e57', 2);

INSERT INTO tree_clusters (id, watering_status, moisture_level, region, address, description, soil_condition, latitude, longitude, geometry)
VALUES 
  (1, 'good', 0.75, 'Mürwik', 'Solitüde Strand', 'Alle Bäume am Strand', 'sandig', 54.820940, 9.489022, ST_SetSRID(ST_MakePoint(54.820940, 9.489022), 4326)),
  (2, 'moderate', 0.5, 'Jürgensby', 'Ulmenstraße', 'Bäume beim Sankt-Jürgen-Platz', 'schluffig', 54.78805731048199, 9.44400186680097, ST_SetSRID(ST_MakePoint(54.78805731048199, 9.44400186680097), 4326));

INSERT INTO sensors (id, status) VALUES (1, 'online');
INSERT INTO sensors (id, status) VALUES (2, 'offline');
INSERT INTO sensors (id, status) VALUES (3, 'unknown');

INSERT INTO trees (tree_cluster_id, sensor_id, age, height_above_sea_level, planting_year, species, tree_number, latitude, longitude, geometry)
VALUES 
  (1, 1, 3, 10.0, 2021, 'Quercus robur', 1, 54.82124518093376, 9.485702120628517, ST_SetSRID(ST_MakePoint(54.82124518093376, 9.485702120628517), 4326)),
  (1, NULL, 2, 11.0, 2022, 'Quercus robur', 1, 54.8215076622281, 9.487153277881877, ST_SetSRID(ST_MakePoint(54.8215076622281, 9.487153277881877), 4326)),
  (1, NULL, 1, 13.3, 2023, 'Quercus robur', 1, 54.82078826498143, 9.489684366114483, ST_SetSRID(ST_MakePoint(54.82078826498143, 9.489684366114483), 4326)),
  (2, 2, 4, 10.0, 2020, 'Quercus robur', 1, 54.78780993841013, 9.444052105200551, ST_SetSRID(ST_MakePoint(54.78780993841013, 9.444052105200551), 4326)),
  (2, NULL, 3, 11.0, 2021, 'Quercus robur', 1, 54.78836553796373, 9.444075995492044, ST_SetSRID(ST_MakePoint(54.78836553796373, 9.444075995492044), 4326)),
  (2, NULL, 2, 13.3, 2022, 'Quercus robur', 1, 54.787768612518455, 9.443996361187065, ST_SetSRID(ST_MakePoint(54.787768612518455, 9.443996361187065), 4326));

INSERT INTO sensor_data (sensor_id, data)
VALUES 
  (1, '{"temperature": 20.0, "humidity": 0.5, "moisture": 0.75}'),
  (1, '{"temperature": 21.0, "humidity": 0.6, "moisture": 0.5}'),
  (1, '{"temperature": 22.0, "humidity": 0.7, "moisture": 0.25}'),
  (2, '{"temperature": 20.0, "humidity": 0.5, "moisture": 0.75}'),
  (2, '{"temperature": 21.0, "humidity": 0.6, "moisture": 0.5}'),
  (2, '{"temperature": 22.0, "humidity": 0.7, "moisture": 0.25}'),
  (3, '{"temperature": 20.0, "humidity": 0.5, "moisture": 0.75}'),
  (3, '{"temperature": 21.0, "humidity": 0.6, "moisture": 0.5}'),
  (3, '{"temperature": 22.0, "humidity": 0.7, "moisture": 0.25}');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM images WHERE url = 'https://avatars.githubusercontent.com/u/165842746?s=96&v=4';
DELETE FROM images WHERE url = 'https://app.dev.green-ecolution.de/api/v1/images/avatar.png';

DELETE FROM user_vehicles WHERE user_id = '95b69b4c-b38b-4394-9520-496879b67791';
DELETE FROM user_vehicles WHERE user_id = 'd2563a8e-a608-4039-8718-25fc3c1d8e57';

DELETE FROM vehicles WHERE number_plate = 'B-1234';
DELETE FROM vehicles WHERE number_plate = 'B-5678';

DELETE FROM trees WHERE tree_cluster_id = 1;
DELETE FROM trees WHERE tree_cluster_id = 2;

DELETE FROM tree_clusters WHERE id = 1;
DELETE FROM tree_clusters WHERE id = 2;

DELETE FROM sensor_data WHERE sensor_id = 1;
DELETE FROM sensor_data WHERE sensor_id = 2;
DELETE FROM sensor_data WHERE sensor_id = 3;

DELETE FROM sensors WHERE id = 1;
DELETE FROM sensors WHERE id = 2;
DELETE FROM sensors WHERE id = 3;

-- +goose StatementEnd
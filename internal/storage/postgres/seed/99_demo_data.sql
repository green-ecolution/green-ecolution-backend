-- +goose Up
-- +goose StatementBegin
INSERT INTO images (url) VALUES ('https://avatars.githubusercontent.com/u/165842746?s=96&v=4');
INSERT INTO images (url, filename, mime_type) VALUES ('https://app.dev.green-ecolution.de/api/v1/images/avatar.png', 'avatar.png', 'image/png');

INSERT INTO vehicles (id, number_plate, description, water_capacity) VALUES (1, 'B-1234', 'Test vehicle 1', 100.0);
INSERT INTO vehicles (id, number_plate, description, water_capacity) VALUES (2, 'B-5678', 'Test vehicle 2', 150.0);
ALTER SEQUENCE vehicles_id_seq RESTART WITH 3;

INSERT INTO user_vehicles (user_id, vehicle_id) VALUES ('95b69b4c-b38b-4394-9520-496879b67791', 1);
INSERT INTO user_vehicles (user_id, vehicle_id) VALUES ('d2563a8e-a608-4039-8718-25fc3c1d8e57', 2);

INSERT INTO tree_clusters (id, name, watering_status, moisture_level, region_id, address, description, soil_condition, latitude, longitude, geometry)
VALUES 
  (1, 'Gruppe: Solitüde Strand', 'good', 0.75, 1, 'Solitüde Strand', 'Alle Bäume am Strand', 'sandig', 54.820940, 9.489022, ST_SetSRID(ST_MakePoint(54.820940, 9.489022), 4326)),
  (2, 'Gruppe: Sankt-Jürgen-Platz', 'moderate', 0.5, 1, 'Ulmenstraße', 'Bäume beim Sankt-Jürgen-Platz', 'schluffig', 54.78805731048199, 9.44400186680097, ST_SetSRID(ST_MakePoint(54.78805731048199, 9.44400186680097), 4326));
ALTER SEQUENCE tree_clusters_id_seq RESTART WITH 3;

INSERT INTO sensors (id, status) VALUES (1, 'online');
INSERT INTO sensors (id, status) VALUES (2, 'offline');
INSERT INTO sensors (id, status) VALUES (3, 'unknown');
INSERT INTO sensors (id, status) VALUES (4, 'online');
ALTER SEQUENCE sensors_id_seq RESTART WITH 5;

INSERT INTO trees (tree_cluster_id, sensor_id, planting_year, species, tree_number, latitude, longitude, geometry, readonly)
VALUES 
  (1, 1, 2021, 'Quercus robur', 1, 54.82124518093376, 9.485702120628517, ST_SetSRID(ST_MakePoint(54.82124518093376, 9.485702120628517), 4326), true),
  (1, NULL, 2022, 'Quercus robur', 1, 54.8215076622281, 9.487153277881877, ST_SetSRID(ST_MakePoint(54.8215076622281, 9.487153277881877), 4326), true),
  (1, NULL, 2023, 'Quercus robur', 1, 54.82078826498143, 9.489684366114483, ST_SetSRID(ST_MakePoint(54.82078826498143, 9.489684366114483), 4326), true),
  (2, 2, 2020, 'Quercus robur', 1, 54.78780993841013, 9.444052105200551, ST_SetSRID(ST_MakePoint(54.78780993841013, 9.444052105200551), 4326), false),
  (2, NULL, 2021, 'Quercus robur', 1, 54.78836553796373, 9.444075995492044, ST_SetSRID(ST_MakePoint(54.78836553796373, 9.444075995492044), 4326), false),
  (2, NULL, 2022, 'Quercus robur', 1, 54.787768612518455, 9.443996361187065, ST_SetSRID(ST_MakePoint(54.787768612518455, 9.443996361187065), 4326), false),
  (2, NULL, 2022, 'Quercus robur', 1, 54.77933725347423, 9.426465409018832, ST_SetSRID(ST_MakePoint(54.77933725347423, 9.426465409018832), 4326), true);

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
DELETE FROM images;
DELETE FROM vehicles;
DELETE FROM trees;
DELETE FROM tree_clusters;
DELETE FROM sensor_data;
DELETE FROM sensors;
-- +goose StatementEnd

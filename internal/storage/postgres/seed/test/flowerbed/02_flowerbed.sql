-- +goose Up
-- +goose StatementBegin
INSERT INTO flowerbeds (id, sensor_id, size, description, number_of_plants, moisture_level, region_id, address, latitude, longitude, geometry) VALUES 
  (1, 2, 20.0, 'Big flowerbed nearby the sea', 10.000, 4.5, 1, '123 Garden streat',  54.776613, 9.454303, ST_GeomFromText('POLYGON((54.776613 9.454303, 54.776713 9.454303, 54.776713 9.454403, 54.776613 9.454403, 54.776613 9.454303))', 4326));
ALTER SEQUENCE flowerbeds_id_seq RESTART WITH 14;

INSERT INTO sensors (id, status) VALUES (1, 'online');
INSERT INTO sensors (id, status) VALUES (2, 'offline');
INSERT INTO sensors (id, status) VALUES (3, 'unknown');
INSERT INTO sensors (id, status) VALUES (4, 'online');
ALTER SEQUENCE sensors_id_seq RESTART WITH 5;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM flowerbeds;
DELETE FROM sensors;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
INSERT INTO flowerbeds (id, sensor_id, size, description, number_of_plants, moisture_level, region_id, address, latitude, longitude, geometry) VALUES 
  (1, 'sensor-2', 20.0, 'Big flowerbed nearby the sea', 10.000, 4.5, 1, '123 Garden street',  54.776613, 9.454303, ST_GeomFromText('POLYGON((54.776613 9.454303, 54.776713 9.454303, 54.776713 9.454403, 54.776613 9.454403, 54.776613 9.454303))', 4326)),
  (2, NULL, 15.0, 'Small flowerbed in the park', 5, 3.2, 3, '456 Park Avenue', 54.776700, 9.454400, ST_GeomFromText('POLYGON((54.776700 9.454400, 54.776800 9.454400, 54.776800 9.454500, 54.776700 9.454500, 54.776700 9.454400))', 4326));
ALTER SEQUENCE flowerbeds_id_seq RESTART WITH 14;

INSERT INTO sensors (id, status, latitude, longitude, geometry)
VALUES
    ('sensor-1', 'online', 54.82124518093376, 9.485702120628517, ST_SetSRID(ST_MakePoint(54.82124518093376, 9.485702120628517), 4326));

INSERT INTO sensors (id, status, latitude, longitude, geometry)
VALUES
    ('sensor-2', 'offline', 54.78780993841013, 9.444052105200551, ST_SetSRID(ST_MakePoint(54.78780993841013, 9.444052105200551), 4326));

INSERT INTO images (id, url, filename, mime_type) VALUES (1, '/test/url/to/image', 'Screenshot', 'png');
ALTER SEQUENCE images_id_seq RESTART WITH 2;

INSERT INTO flowerbed_images (flowerbed_id, image_id) VALUES (1, 1);
ALTER SEQUENCE flowerbeds_id_seq RESTART WITH 14;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM flowerbeds;
DELETE FROM sensors;
DELETE FROM images;
DELETE FROM flowerbed_images;
-- +goose StatementEnd

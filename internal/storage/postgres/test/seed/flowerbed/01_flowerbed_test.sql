-- +goose Up
-- +goose StatementBegin

INSERT INTO flowerbeds (id, sensor_id, size, description, number_of_plants, moisture_level, region, address, archived, latitude, longitude)
VALUES
  (1, 1, 10.0, 'Blumenbeet am Strand', 5, 0.75, 'Mürwik', 'Solitüde Strand', false, 54.820940, 9.489022),
  (2, 2, 11.0, 'Blumenbeet beim Sankt-Jürgen-Platz', 5, 0.5, 'Jürgensby', 'Ulmenstraße', false, 54.78805731048199, 9.44400186680097);
ALTER SEQUENCE flowerbeds_id_seq RESTART WITH 3;

INSERT INTO sensors (id, status) VALUES (1, 'online');
INSERT INTO sensors (id, status) VALUES (2, 'offline');
ALTER SEQUENCE sensors_id_seq RESTART WITH 3;

INSERT INTO images (id, url) VALUES (1, 'https://avatars.githubusercontent.com/u/165842746?s=96&v=4');
INSERT INTO images (id, url, filename, mime_type) VALUES (2, 'https://app.dev.green-ecolution.de/api/v1/images/flowerbed.png', 'flowerbed.png', 'image/png');
ALTER SEQUENCE images_id_seq RESTART WITH 3;

INSERT INTO flowerbed_images (flowerbed_id, image_id) VALUES (1, 1);
INSERT INTO flowerbed_images (flowerbed_id, image_id) VALUES (1, 2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM flowerbed_images WHERE flowerbed_id = 1;
DELETE FROM flowerbed_images WHERE flowerbed_id = 2;

DELETE FROM sensors WHERE id = 1;
DELETE FROM sensors WHERE id = 2;

DELETE FROM images WHERE id = 1;
DELETE FROM images WHERE id = 2;

DELETE FROM flowerbeds WHERE id = 1;
DELETE FROM flowerbeds WHERE id = 2;
-- +goose StatementEnd

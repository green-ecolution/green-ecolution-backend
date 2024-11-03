-- +goose Up
-- +goose StatementBegin
INSERT INTO tree_clusters (id, name, watering_status, moisture_level, region_id, address, description, soil_condition, latitude, longitude, geometry)
VALUES
    (1, 'Solitüde Strand', 'good', 0.75, 1, 'Solitüde Strand', 'Alle Bäume am Strand', 'sandig', 54.820940, 9.489022, ST_SetSRID(ST_MakePoint(54.820940, 9.489022), 4326)),
    (2, 'Sankt-Jürgen-Platz', 'moderate', 0.5, 1, 'Ulmenstraße', 'Bäume beim Sankt-Jürgen-Platz', 'schluffig', 54.78805731048199, 9.44400186680097, ST_SetSRID(ST_MakePoint(54.78805731048199, 9.44400186680097), 4326)),
    (3, 'Flensburger Stadion', 'unknown', 0.7, 1, 'Flensburger Stadion', 'Alle Bäume in der Gegend des Stadions in Mürwik', 'schluffig', 54.802163, 9.446398, ST_SetSRID(ST_MakePoint(54.802163, 9.446398), 4326)),
    (4, 'Campus Hochschule', 'bad', 0.1, 4, 'Thomas-Finke Straße', 'Gruppe ist besonders anfällig', 'schluffig', 54.776613, 9.454303, ST_SetSRID(ST_MakePoint(54.776613, 9.454303), 4326)),
    (5, 'Mathildenstraße', 'moderate', 0.4, 10, 'Mathildenstraße', 'Sehr enge Straße und dadurch schlecht zu bewässern.', 'schluffig', 54.782402, 9.424270, ST_SetSRID(ST_MakePoint(54.782402, 9.424270), 4326)),
    (6, 'Nordstadt', 'good', 0.6, 13, 'Apenrader Straße', 'Guter Baumbestand mit großen Kronen.', 'sandig', 54.807162, 9.423138, ST_SetSRID(ST_MakePoint(54.807162, 9.423138), 4326)),
    (7, 'TSB Neustadt', 'good', 0.75, 13, 'Ecknerstraße', 'Kleiner Baumbestand.', 'sandig', 54.797162, 9.419620, ST_SetSRID(ST_MakePoint(54.797162, 9.419620), 4326)),
    (8, 'Gewerbegebiet Süd', 'bad', 0.1, NULL, 'Schleswiger Straße', 'Sehr viel versiegelter Boden.', 'sandig', 54.768115, 9.435285, ST_SetSRID(ST_MakePoint(54.768115, 9.435285), 4326));
ALTER SEQUENCE tree_clusters_id_seq RESTART WITH 9;

INSERT INTO sensors (id, status) VALUES (1, 'online');
INSERT INTO sensors (id, status) VALUES (2, 'offline');
INSERT INTO sensors (id, status) VALUES (3, 'unknown');
INSERT INTO sensors (id, status) VALUES (4, 'online');
ALTER SEQUENCE sensors_id_seq RESTART WITH 5;

INSERT INTO images (id, url, filename, mime_type)
VALUES
    (1, 'example/image1.jpg', 'image1.jpg', 'image/jpeg'),
    (2, 'example/image2.jpg', 'image2.jpg', 'image/jpeg');

ALTER SEQUENCE images_id_seq RESTART WITH 3;

INSERT INTO trees (id, tree_cluster_id, sensor_id, planting_year, species, tree_number, latitude, longitude, geometry, readonly, watering_status, description)
VALUES
    (1, 1, 1, 2021, 'Quercus robur', 1005, 54.82124518093376, 9.485702120628517, ST_SetSRID(ST_MakePoint(54.82124518093376, 9.485702120628517), 4326), true, 'unknown', 'Sample description 1'),
    (2, 1, NULL, 2022, 'Quercus robur', 1006, 54.8215076622281, 9.487153277881877, ST_SetSRID(ST_MakePoint(54.8215076622281, 9.487153277881877), 4326), true, 'good', 'Sample description 2'),
    (3, 2, 2, 2023, 'Betula pendula', 1010, 54.78780993841013, 9.444052105200551, ST_SetSRID(ST_MakePoint(54.78780993841013, 9.444052105200551), 4326), false, 'bad', 'Sample description 3');

ALTER SEQUENCE trees_id_seq RESTART WITH 4;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM trees;
DELETE FROM tree_clusters;
DELETE FROM sensor_data;
-- +goose StatementEnd

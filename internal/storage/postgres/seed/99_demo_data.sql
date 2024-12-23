-- +goose Up
-- +goose StatementBegin
INSERT INTO images (url) VALUES ('https://avatars.githubusercontent.com/u/165842746?s=96&v=4');
INSERT INTO images (url, filename, mime_type) VALUES ('https://app.dev.green-ecolution.de/api/v1/images/avatar.png', 'avatar.png', 'image/png');

INSERT INTO vehicles (id, number_plate, description, water_capacity, type, status, driving_license, model, width, height, length) VALUES (1, 'B-1234', 'Test vehicle 1', 100.0, 'trailer', 'active', 'BE', 'LK1615/17 - Conrad - MAN TGE 3.180', 2.0, 1.5, 2.0);
INSERT INTO vehicles (id, number_plate, description, water_capacity, type, status, driving_license, model, width, height, length) VALUES (2, 'B-5678', 'Test vehicle 2', 150.0, 'transporter', 'unknown', 'C', 'Actros L Mercedes Benz', 2.4, 2.1, 5.0);
ALTER SEQUENCE vehicles_id_seq RESTART WITH 3;

INSERT INTO tree_clusters (id, name, watering_status, moisture_level, region_id, address, description, soil_condition, latitude, longitude, geometry)
VALUES 
  (1, 'Solitüde Strand', 'good', 0.75, 1, 'Solitüde Strand', 'Alle Bäume am Strand', 'sandig', 54.820940, 9.489022, ST_SetSRID(ST_MakePoint(54.820940, 9.489022), 4326)),
  (2, 'Sankt-Jürgen-Platz', 'moderate', 0.5, 1, 'Ulmenstraße', 'Bäume beim Sankt-Jürgen-Platz', 'schluffig', 54.78805731048199, 9.44400186680097, ST_SetSRID(ST_MakePoint(54.78805731048199, 9.44400186680097), 4326)),
  (3, 'Flensburger Stadion', 'unknown', 0.7, 1, 'Flensburger Stadion', 'Alle Bäume in der Gegend des Stadions in Mürwik', 'schluffig', 54.802163, 9.446398, ST_SetSRID(ST_MakePoint(54.802163, 9.446398), 4326)),
  (4, 'Campus Hochschule', 'bad', 0.1, 4, 'Thomas-Finke Straße', 'Gruppe ist besonders anfällig', 'schluffig', 54.776613, 9.454303, ST_SetSRID(ST_MakePoint(54.776613, 9.454303), 4326)),
  (5, 'Mathildenstraße', 'moderate', 0.4, 10, 'Mathildenstraße', 'Sehr enge Straße und dadurch schlecht zu bewässern.', 'schluffig', 54.782402, 9.424270, ST_SetSRID(ST_MakePoint(54.782402, 9.424270), 4326)),
  (6, 'Nordstadt', 'good', 0.6, 13, 'Apenrader Straße', 'Guter Baumbestand mit großen Kronen.', 'sandig', 54.807162, 9.423138, ST_SetSRID(ST_MakePoint(54.807162, 9.423138), 4326)),
  (7, 'TSB Neustadt', 'good', 0.75, 13, 'Ecknerstraße', 'Kleiner Baumbestand.', 'sandig', 54.797162, 9.419620, ST_SetSRID(ST_MakePoint(54.797162, 9.419620), 4326)),
  (8, 'Gewerbegebiet Süd', 'bad', 0.1, 13, 'Schleswiger Straße', 'Sehr viel versiegelter Boden.', 'sandig', 54.768115, 9.435285, ST_SetSRID(ST_MakePoint(54.768115, 9.435285), 4326));   
ALTER SEQUENCE tree_clusters_id_seq RESTART WITH 9;

INSERT INTO sensors (id, status, latitude, longitude, geometry)
VALUES
    ('sensor-1', 'online', 54.82124518093376, 9.485702120628517, ST_SetSRID(ST_MakePoint(54.82124518093376, 9.485702120628517), 4326));

INSERT INTO sensors (id, status, latitude, longitude, geometry)
VALUES
    ('sensor-2', 'offline', 54.78780993841013, 9.444052105200551, ST_SetSRID(ST_MakePoint(54.78780993841013, 9.444052105200551), 4326));

INSERT INTO sensors (id, status, latitude, longitude, geometry)
VALUES
    ('sensor-3', 'unknown', 54.779832339352, 9.424575519708398, ST_SetSRID(ST_MakePoint(54.779832339352, 9.424575519708398), 4326));

INSERT INTO sensors (id, status, latitude, longitude, geometry)
VALUES
    ('sensor-4', 'online', 54.82042158396298, 9.488535474607701, ST_SetSRID(ST_MakePoint(54.82042158396298, 9.488535474607701), 4326));

INSERT INTO sensors (id, status, latitude, longitude, geometry)
VALUES
    ('tree-sensor', 'online', 54.774932, 9.450000, ST_SetSRID(ST_MakePoint(54.774932, 9.450000), 4326));

INSERT INTO trees (tree_cluster_id, sensor_id, planting_year, species, number, latitude, longitude, geometry, readonly, watering_status, description)
VALUES 
  (1, 'sensor-1', 2021, 'Quercus robur', 1005, 54.82124518093376, 9.485702120628517, ST_SetSRID(ST_MakePoint(54.82124518093376, 9.485702120628517), 4326), true, 'unknown', 'Dieser Baum wurde im August das lezte mal gestuzt'),
  (1, NULL, 2022, 'Quercus robur', 1006, 54.8215076622281, 9.487153277881877, ST_SetSRID(ST_MakePoint(54.8215076622281, 9.487153277881877), 4326), true, 'good', ''),
  (1, NULL, 2023, 'Quercus robur', 1007, 54.82078826498143, 9.489684366114483, ST_SetSRID(ST_MakePoint(54.82078826498143, 9.489684366114483), 4326), true, 'moderate', ''),
  (2, 'sensor-2', 2020, 'Quercus robur', 1008, 54.78780993841013, 9.444052105200551, ST_SetSRID(ST_MakePoint(54.78780993841013, 9.444052105200551), 4326), false, 'bad', ''),
  (2, NULL, 2021, 'Quercus robur', 1009, 54.78836553796373, 9.444075995492044, ST_SetSRID(ST_MakePoint(54.78836553796373, 9.444075995492044), 4326), false, 'unknown', ''),
  (2, NULL, 2022, 'Quercus robur', 1010, 54.787768612518455, 9.443996361187065, ST_SetSRID(ST_MakePoint(54.787768612518455, 9.443996361187065), 4326), false, 'good', ''),
  (NULL, NULL, 2022, 'Quercus robur', 1011, 54.77933725347423, 9.426465409018832, ST_SetSRID(ST_MakePoint(54.77933725347423, 9.426465409018832), 4326), true, 'moderate', ''),
  (NULL, NULL, 2024, 'Tilia cordata', 1012, 54.782528, 9.424376, ST_SetSRID(ST_MakePoint(54.782528, 9.424376), 4326), true, 'bad', 'Dieser Baum wurde im August das lezte mal gestuzt'),
  (NULL, NULL, 2024, 'Tilia x vulgaris', 1013, 54.782667, 9.424816, ST_SetSRID(ST_MakePoint(54.782667, 9.424816), 4326), true, 'bad', ''),
  (NULL, NULL, 2024, 'Tilia x vulgaris', 1014, 54.782925, 9.425931, ST_SetSRID(ST_MakePoint(54.782925, 9.425931), 4326), true, 'good', ''),
  (NULL, NULL, 2023, 'Tilia x vulgaris', 1015, 54.783739, 9.426823, ST_SetSRID(ST_MakePoint(54.783739, 9.426823), 4326), true, 'unknown', ''),
  (NULL, NULL, 2023, 'Tilia cordata', 1016, 54.783739, 9.426823, ST_SetSRID(ST_MakePoint(54.783739, 9.426823), 4326), true, 'good', ''),
  (NULL, NULL, 2022, 'Carpinus betulus', 1017, 54.785981, 9.430668, ST_SetSRID(ST_MakePoint(54.785981, 9.430668), 4326), true, 'unknown', ''),
  (NULL, NULL, 2022, 'Carpinus betulus', 1018, 54.786269, 9.431758, ST_SetSRID(ST_MakePoint(54.786269, 9.431758), 4326), true, 'unknown', ''),
  (NULL, NULL, 2022, 'Carpinus betulus', 1019, 54.787339, 9.431701, ST_SetSRID(ST_MakePoint(54.787339, 9.431701), 4326), true, 'unknown', ''),
  (NULL, NULL, 2021, 'Carpinus betulus', 1020, 54.786656, 9.432243, ST_SetSRID(ST_MakePoint(54.786656, 9.432243), 4326), true, 'good', 'Dieser Baum wurde im August das lezte mal gestuzt'),
  (5, NULL, 2021, 'Fraxinus ornus Obelisk', 1021, 54.782630, 9.423792, ST_SetSRID(ST_MakePoint(54.782630, 9.423792), 4326), true, 'good', ''),
  (5, NULL, 2021, 'Fraxinus ornus Obelisk', 1022, 54.782463, 9.423727, ST_SetSRID(ST_MakePoint(54.782463, 9.423727), 4326), true, 'unknown', ''),
  (5, NULL, 2021, 'Fraxinus ornus Obelisk', 1023, 54.782296, 9.424178, ST_SetSRID(ST_MakePoint(54.782296, 9.424178), 4326), true, 'unknown', 'Dieser Baum wurde im August das lezte mal gestuzt'),
  (5, NULL, 2022, 'Fraxinus ornus Obelisk', 1024, 54.782043, 9.424188, ST_SetSRID(ST_MakePoint(54.782043, 9.424188), 4326), true, 'moderate', ''),
  (5, NULL, 2022, 'Fraxinus ornus Obelisk', 1025, 54.781753, 9.424936, ST_SetSRID(ST_MakePoint(54.781753, 9.424936), 4326), true, 'moderate', ''),
  (4, NULL, 2022, 'Tilia intermedia', 1026, 54.776145, 9.449785, ST_SetSRID(ST_MakePoint(54.776145, 9.449785), 4326), true, 'moderate', ''),
  (4, NULL, 2021, 'Tilia intermedia', 1027, 54.776120, 9.450891, ST_SetSRID(ST_MakePoint(54.776120, 9.450891), 4326), true, 'unknown', 'Dieser Baum wurde im August das lezte mal gestuzt'),
  (4, NULL, 2021, 'Tilia intermedia', 1028, 54.776058, 9.450311, ST_SetSRID(ST_MakePoint(54.776058, 9.450311), 4326), true, 'bad', ''),
  (4, NULL, 2021, 'Tilia intermedia', 1029, 54.775709, 9.447762, ST_SetSRID(ST_MakePoint(54.775709, 9.447762), 4326), true, 'unknown', 'Dieser Baum wurde im August das lezte mal gestuzt'),
  (NULL, NULL, 2023, 'Alnus glutinosa', 1030, 54.792472, 9.452773, ST_SetSRID(ST_MakePoint(54.792472, 9.452773), 4326), true, 'unknown', ''),
  (NULL, NULL, 2023, 'Alnus glutinosa', 1031, 54.792782, 9.453795, ST_SetSRID(ST_MakePoint(54.792782, 9.453795), 4326), true, 'unknown', 'Dieser Baum wurde im August das lezte mal gestuzt'),
  (NULL, NULL, 2023, 'Alnus glutinosa', 1032, 54.792837, 9.454880, ST_SetSRID(ST_MakePoint(54.792837, 9.454880), 4326), true, 'good', ''),
  (NULL, NULL, 2023, 'Alnus glutinosa', 1033, 54.792435, 9.455545, ST_SetSRID(ST_MakePoint(54.792435, 9.455545), 4326), true, 'good', ''),
  (3, NULL, 2023, 'Betula pendula', 1034, 54.801718, 9.444797, ST_SetSRID(ST_MakePoint(54.801718, 9.444797), 4326), true, 'good', ''),
  (3, NULL, 2023, 'Betula pendula', 1035, 54.800797, 9.444271, ST_SetSRID(ST_MakePoint(54.800797, 9.444271), 4326), true, 'unknown', 'Dieser Baum wurde im August das lezte mal gestuzt'),
  (3, NULL, 2023, 'Betula pendula', 1036, 54.801539, 9.446741, ST_SetSRID(ST_MakePoint(54.801539, 9.446741), 4326), true, 'unknown', ''),
  (3, NULL, 2023, 'Betula pendula', 1037, 54.799796, 9.443927, ST_SetSRID(ST_MakePoint(54.799796, 9.443927), 4326), true, 'unknown', ''),
  (3, NULL, 2023, 'Betula pendula', 1038, 54.804052, 9.447900, ST_SetSRID(ST_MakePoint(54.804052, 9.447900), 4326), true, 'unknown', ''),
  (6, NULL, 2023, 'Acer platanoides Schwedleri', 1039, 54.806287, 9.423469, ST_SetSRID(ST_MakePoint(54.806287, 9.423469), 4326), true, 'unknown', ''),
  (6, NULL, 2023, 'Acer platanoides Schwedleri', 1040, 54.807212, 9.422752, ST_SetSRID(ST_MakePoint(54.807212, 9.422752), 4326), true, 'unknown', ''),
  (6, NULL, 2023, 'Acer platanoides Schwedleri', 1041, 54.806606, 9.422773, ST_SetSRID(ST_MakePoint(54.806606, 9.422773), 4326), true, 'good', 'Dieser Baum wurde im August das lezte mal gestuzt'),
  (6, NULL, 2023, 'Acer platanoides Schwedleri', 1042, 54.807787, 9.422354, ST_SetSRID(ST_MakePoint(54.807787, 9.422354), 4326), true, 'good', 'Dieser Baum wurde im August das lezte mal gestuzt'),
  (7, NULL, 2022, 'Acer platanoides Schwedleri', 1043, 54.796916, 9.421332, ST_SetSRID(ST_MakePoint(54.796916, 9.421332), 4326), true, 'good', ''),
  (7, NULL, 2022, 'Acer platanoides Schwedleri', 1044, 54.797330, 9.419002, ST_SetSRID(ST_MakePoint(54.797330, 9.419002), 4326), true, 'unknown', ''),
  (7, NULL, 2022, 'Acer platanoides Schwedleri', 1045, 54.797114, 9.417843, ST_SetSRID(ST_MakePoint(54.797114, 9.417843), 4326), true, 'unknown', ''),
  (8, NULL, 2022, 'Sorbus x thuringiaca', 1046, 54.767936, 9.435316, ST_SetSRID(ST_MakePoint(54.767936, 9.435316), 4326), true, 'unknown', ''),
  (8, NULL, 2022, 'Sorbus x thuringiaca', 1047, 54.767275, 9.435024, ST_SetSRID(ST_MakePoint(54.767275, 9.435024), 4326), true, 'unknown', 'Dieser Baum wurde im August das lezte mal gestuzt'),
  (8, NULL, 2022, 'Sorbus x thuringiaca', 1048, 54.766991, 9.435672, ST_SetSRID(ST_MakePoint(54.766991, 9.435672), 4326), true, 'unknown', ''),
  (8, NULL, 2022, 'Sorbus x thuringiaca', 1049, 54.767972, 9.435373, ST_SetSRID(ST_MakePoint(54.767972, 9.435373), 4326), true, 'bad', ''),
  (8, NULL, 2022, 'Sorbus x thuringiaca', 1050, 54.767019, 9.435321, ST_SetSRID(ST_MakePoint(54.767019, 9.435321), 4326), true, 'moderate', 'Dieser Baum wurde im August das lezte mal gestuzt'), 
  (NULL, NULL, 2022, 'Populus cf. suaveolens', 1051, 54.769030, 9.429936, ST_SetSRID(ST_MakePoint(54.769030, 9.429936), 4326), true, 'unknown', ''),
  (NULL, NULL, 2023, 'Populus cf. suaveolens', 1052, 54.768518, 9.439067, ST_SetSRID(ST_MakePoint(54.768518, 9.439067), 4326), true, 'unknown', ''),
  (NULL, NULL, 2021, 'Populus cf. suaveolens', 1053, 54.775237, 9.441981, ST_SetSRID(ST_MakePoint(54.775237, 9.441981), 4326), true, 'unknown', ''),
  (NULL, NULL, 2022, 'Populus cf. suaveolens', 1054, 54.780192, 9.459607, ST_SetSRID(ST_MakePoint(54.780192, 9.459607), 4326), true, 'unknown', 'Dieser Baum wurde im August das lezte mal gestuzt'),
  (NULL, NULL, 2021, 'Populus cf. suaveolens', 1055, 54.785043, 9.418210, ST_SetSRID(ST_MakePoint(54.785043, 9.418210), 4326), true, 'unknown', ''),
  (NULL, NULL, 2022, 'Fraxinus excelsior', 1056, 54.779697, 9.440026, ST_SetSRID(ST_MakePoint(54.779697, 9.440026), 4326), true, 'unknown', ''),
  (NULL, NULL, 2020, 'Fraxinus excelsior', 1057, 54.785147, 9.438903, ST_SetSRID(ST_MakePoint(54.785147, 9.438903), 4326), true, 'unknown', ''),
  (NULL, NULL, 2020, 'Fraxinus excelsior', 1058, 54.788205, 9.454699, ST_SetSRID(ST_MakePoint(54.788205, 9.454699), 4326), true, 'unknown', 'Dieser Baum wurde im August das lezte mal gestuzt'),
  (NULL, NULL, 2020, 'Fraxinus excelsior', 1059, 54.804054, 9.469544, ST_SetSRID(ST_MakePoint(54.804054, 9.469544), 4326), true, 'unknown', 'Dieser Baum wurde im August das lezte mal gestuzt'),
  (NULL, NULL, 2024, 'Acer pseudoplatanus', 1060, 54.813655, 9.477633, ST_SetSRID(ST_MakePoint(54.813655, 9.477633), 4326), true, 'unknown', 'Dieser Baum wurde im August das lezte mal gestuzt'),
  (NULL, NULL, 2024, 'Acer pseudoplatanus', 1061, 54.811001, 9.484132, ST_SetSRID(ST_MakePoint(54.811001, 9.484132), 4326), true, 'unknown', 'Dieser Baum wurde im August das lezte mal gestuzt'),
  (NULL, NULL, 2024, 'Acer pseudoplatanus', 1062, 54.790366, 9.472744, ST_SetSRID(ST_MakePoint(54.790366, 9.472744), 4326), true, 'unknown', '');

INSERT INTO sensor_data (sensor_id, data)
VALUES
    ('sensor-1', '{
        "temperature": 20.0,
        "humidity": 0.5,
        "battery": 34,
        "watermarks": [
            {"resistance": 23, "centibar": 38, "depth": 30},
            {"resistance": 24, "centibar": 40, "depth": 60},
            {"resistance": 25, "centibar": 42, "depth": 90}
        ]
    }'),
    ('sensor-1', '{
        "temperature": 21.5,
        "humidity": 0.6,
        "battery": 36,
        "watermarks": [
            {"resistance": 26, "centibar": 45, "depth": 30},
            {"resistance": 28, "centibar": 47, "depth": 60},
            {"resistance": 30, "centibar": 50, "depth": 90}
        ]
    }'),
    ('sensor-1', '{
        "temperature": 22.0,
        "humidity": 0.7,
        "battery": 32,
        "watermarks": [
            {"resistance": 29, "centibar": 55, "depth": 30},
            {"resistance": 31, "centibar": 58, "depth": 60},
            {"resistance": 33, "centibar": 60, "depth": 90}
        ]
    }'),
    ('sensor-2', '{
        "temperature": 20.0,
        "humidity": 0.5,
        "battery": 40,
        "watermarks": [
            {"resistance": 22, "centibar": 35, "depth": 30},
            {"resistance": 23, "centibar": 38, "depth": 60},
            {"resistance": 24, "centibar": 40, "depth": 90}
        ]
    }'),
    ('sensor-2', '{
        "temperature": 21.0,
        "humidity": 0.6,
        "battery": 39,
        "watermarks": [
            {"resistance": 25, "centibar": 43, "depth": 30},
            {"resistance": 27, "centibar": 46, "depth": 60},
            {"resistance": 28, "centibar": 49, "depth": 90}
        ]
    }'),
    ('sensor-2', '{
        "temperature": 22.0,
        "humidity": 0.7,
        "battery": 38,
        "watermarks": [
            {"resistance": 29, "centibar": 52, "depth": 30},
            {"resistance": 30, "centibar": 55, "depth": 60},
            {"resistance": 32, "centibar": 58, "depth": 90}
        ]
    }'),
    ('sensor-3', '{
        "temperature": 19.5,
        "humidity": 0.4,
        "battery": 37,
        "watermarks": [
            {"resistance": 20, "centibar": 32, "depth": 30},
            {"resistance": 21, "centibar": 34, "depth": 60},
            {"resistance": 22, "centibar": 36, "depth": 90}
        ]
    }'),
    ('sensor-3', '{
        "temperature": 20.5,
        "humidity": 0.5,
        "battery": 36,
        "watermarks": [
            {"resistance": 23, "centibar": 38, "depth": 30},
            {"resistance": 24, "centibar": 40, "depth": 60},
            {"resistance": 25, "centibar": 42, "depth": 90}
        ]
    }'),
    ('sensor-3', '{
        "temperature": 21.5,
        "humidity": 0.6,
        "battery": 35,
        "watermarks": [
            {"resistance": 26, "centibar": 45, "depth": 30},
            {"resistance": 28, "centibar": 47, "depth": 60},
            {"resistance": 30, "centibar": 50, "depth": 90}
        ]
    }');


INSERT INTO watering_plans (id, date, description, status, distance, total_water_required, cancellation_note)
VALUES 
  (1, '2024-09-22', 'New watering plan for the west side of the city', 'planned', 63.0, 720.0, ''),
  (2, '2024-08-03', 'New watering plan for the east side of the city', 'active', 63.0, 600.0, ''),
  (3, '2024-06-12', 'Very important watering plan due to no rainfall', 'finished', 63.0, 1320.0, ''),
  (4, '2024-06-10', 'New watering plan for the south side of the city', 'not competed', 63.0, 600.0, ''),
  (5, '2024-06-04', 'Canceled due to flood', 'canceled', 63.0, 600.0, 'The watering plan was cancelled due to various reasons.');
ALTER SEQUENCE watering_plans_id_seq RESTART WITH 6;

INSERT INTO vehicle_watering_plans (vehicle_id, watering_plan_id) 
VALUES 
  (1, 1),
  (2, 1),
  (2, 2),
  (2, 3),
  (2, 4),
  (2, 5);

INSERT INTO tree_cluster_watering_plans (tree_cluster_id, watering_plan_id)
VALUES 
  (1, 1),
  (2, 1),
  (3, 2),
  (1, 3),
  (2, 3),
  (3, 3),
  (3, 4),
  (3, 5);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM images;
DELETE FROM vehicles;
DELETE FROM trees;
DELETE FROM tree_clusters;
DELETE FROM sensor_data;
DELETE FROM sensors;
DELETE FROM watering_plans;
DELETE FROM vehicle_watering_plans;
DELETE FROM tree_cluster_watering_plans;
-- +goose StatementEnd

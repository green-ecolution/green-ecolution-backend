-- +goose Up
ALTER TABLE sensor_data DROP CONSTRAINT IF EXISTS sensor_data_sensor_id_fkey;
ALTER TABLE trees DROP CONSTRAINT IF EXISTS trees_sensor_id_fkey;
ALTER TABLE flowerbeds DROP CONSTRAINT IF EXISTS flowerbeds_sensor_id_fkey;

ALTER TABLE sensor_data ALTER COLUMN sensor_id TYPE VARCHAR USING sensor_id::VARCHAR;
ALTER TABLE trees ALTER COLUMN sensor_id TYPE VARCHAR USING sensor_id::VARCHAR;
ALTER TABLE flowerbeds ALTER COLUMN sensor_id TYPE VARCHAR USING sensor_id::VARCHAR;

ALTER TABLE sensors ALTER COLUMN id TYPE VARCHAR USING id::VARCHAR;

ALTER TABLE sensor_data ADD CONSTRAINT sensor_data_sensor_id_fkey FOREIGN KEY (sensor_id) REFERENCES sensors(id);
ALTER TABLE trees ADD CONSTRAINT trees_sensor_id_fkey FOREIGN KEY (sensor_id) REFERENCES sensors(id);
ALTER TABLE flowerbeds ADD CONSTRAINT flowerbeds_sensor_id_fkey FOREIGN KEY (sensor_id) REFERENCES sensors(id);

ALTER TABLE sensors
    ADD COLUMN latitude FLOAT NOT NULL,
    ADD COLUMN longitude FLOAT NOT NULL,
    ADD COLUMN geometry GEOMETRY(Point, 4326);
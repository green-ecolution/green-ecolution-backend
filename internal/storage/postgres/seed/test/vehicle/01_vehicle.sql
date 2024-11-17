-- +goose Up
-- +goose StatementBegin
INSERT INTO vehicles (id, number_plate, description, water_capacity, type, status) VALUES (1, 'B-1234', 'Test vehicle 1', 100.0, 'trailer', 'active');
INSERT INTO vehicles (id, number_plate, description, water_capacity, type, status) VALUES (2, 'B-5678', 'Test vehicle 2', 150.0, 'transporter', 'unknown');
ALTER SEQUENCE vehicles_id_seq RESTART WITH 3;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM vehicles;
ALTER SEQUENCE vehicles_id_seq RESTART WITH 1;
-- +goose StatementEnd

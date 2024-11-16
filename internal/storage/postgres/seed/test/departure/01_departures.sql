-- +goose Up
-- +goose StatementBegin
INSERT INTO depatures (id, name, description, latitude, longitude) VALUES (
  (1, "Kielseng", "Nebenzentrale", 54.805122, 9.448282),
  (2, "Schleswiger Straße", "Hauptzentrale", 54.768270, 9.436714)
);

ALTER SEQUENCE departures_id_seq RESTART WITH 3;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM depatures;
ALTER SEQUENCE departures_id_seq RESTART WITH 1;
-- +goose StatementEnd

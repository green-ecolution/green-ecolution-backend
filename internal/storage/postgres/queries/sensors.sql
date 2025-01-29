-- name: GetAllSensors :many
SELECT * FROM sensors ORDER BY id;

-- name: GetAllSensorsByProvider :many
SELECT * FROM sensors WHERE provider = $1 ORDER BY id;

-- name: GetSensorByID :one
SELECT * FROM sensors WHERE id = $1;

-- name: GetSensorByStatus :many
SELECT * FROM sensors WHERE status = $1;

-- name: GetLatestSensorDataByID :one
SELECT *
FROM sensor_data
WHERE sensor_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: CreateSensor :one
INSERT INTO sensors (
    id, status, latitude, longitude
) VALUES (
  $1, $2, $3, $4
) RETURNING id;

-- name: UpdateSensor :exec
UPDATE sensors SET
  status = $2
WHERE id = $1;

-- name: SetSensorLocation :exec
UPDATE sensors SET
    latitude = $2,
    longitude = $3,
    geometry = ST_SetSRID(ST_MakePoint($2, $3), 4326)
WHERE id = $1;

-- name: InsertSensorData :exec
INSERT INTO sensor_data (
  sensor_id, data
) VALUES (
  $1, $2
) RETURNING id;

-- name: DeleteSensor :exec
DELETE FROM sensors WHERE id = $1;

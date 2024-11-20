-- name: GetAllSensors :many
SELECT * FROM sensors;

-- name: GetSensorByID :one
SELECT * FROM sensors WHERE id = $1;

-- name: GetSensorByStatus :many
SELECT * FROM sensors WHERE status = $1;

-- name: GetSensorDataBySensorID :many
SELECT * FROM sensor_data WHERE sensor_id = $1;

-- name: CreateSensor :one
INSERT INTO sensors (
    id,
  status
) VALUES (
  $1,
$2
) RETURNING id;

-- name: UpdateSensor :exec
UPDATE sensors SET
  status = $2
WHERE id = $1;

-- name: InsertSensorData :exec
INSERT INTO sensor_data (
  sensor_id, data 
) VALUES (
  $1, $2
) RETURNING id;

-- name: DeleteSensor :exec
DELETE FROM sensors WHERE id = $1;

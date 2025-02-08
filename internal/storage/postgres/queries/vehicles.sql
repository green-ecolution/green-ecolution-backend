-- name: GetAllVehicles :many
SELECT * FROM vehicles 
ORDER BY water_capacity DESC
LIMIT $1 OFFSET $2;

-- name: GetAllVehiclesCount :one
SELECT COUNT(*) FROM vehicles;

-- name: GetAllVehiclesByProvider :many
SELECT * FROM vehicles WHERE provider = $1 ORDER BY water_capacity DESC;

-- name: GetAllVehiclesByType :many
SELECT * FROM vehicles 
WHERE type = $1 
ORDER BY water_capacity DESC
LIMIT $2 OFFSET $3;

-- name: GetAllVehiclesByTypeCount :one
SELECT COUNT(*) FROM vehicles WHERE type = $1;

-- name: GetVehicleByID :one
SELECT * FROM vehicles WHERE id = $1;

-- name: GetVehicleByPlate :one
SELECT * FROM vehicles WHERE number_plate = $1;

-- name: CreateVehicle :one
INSERT INTO vehicles (
  number_plate, 
  description, 
  water_capacity, 
  type, 
  status, 
  model, 
  driving_license, 
  height, 
  length, 
  width,
  weight,
  provider,
  additional_informations
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
) RETURNING id;

-- name: UpdateVehicle :exec
UPDATE vehicles SET
  number_plate = $2,
  description = $3,
  water_capacity = $4,
  type = $5,
  status = $6,
  model = $7,
  driving_license = $8,
  height = $9,
  length = $10,
  width = $11,
  weight = $12,
  provider = $13,
  additional_informations = $14
WHERE id = $1;

-- name: DeleteVehicle :one
DELETE FROM vehicles WHERE id = $1 RETURNING id;

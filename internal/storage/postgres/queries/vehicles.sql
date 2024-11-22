-- name: GetAllVehicles :many
SELECT * FROM vehicles;

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
  driver_licence, 
  height, 
  length, 
  width
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING id;

-- name: UpdateVehicle :exec
UPDATE vehicles SET
  number_plate = $2,
  description = $3,
  water_capacity = $4,
  type = $5,
  status = $6,
  model = $7,
  driver_licence = $8,
  height = $9,
  length = $10,
  width = $11
WHERE id = $1;

-- name: DeleteVehicle :one
DELETE FROM vehicles WHERE id = $1 RETURNING id;

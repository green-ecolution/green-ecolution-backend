-- name: GetAllVehicles :many
SELECT * FROM vehicles;

-- name: GetVehicleByID :one
SELECT * FROM vehicles WHERE id = $1;

-- name: GetVehicleByPlate :one
SELECT * FROM vehicles WHERE number_plate = $1;

-- name: CreateVehicle :one
INSERT INTO vehicles (
  number_plate, description, water_capacity
) VALUES (
  $1, $2, $3
) RETURNING id;

-- name: UpdateVehicle :exec
UPDATE vehicles SET
  number_plate = $2,
  description = $3,
  water_capacity = $4
WHERE id = $1;

-- name: DeleteVehicle :exec
DELETE FROM vehicles WHERE id = $1;

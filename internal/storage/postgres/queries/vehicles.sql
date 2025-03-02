-- name: GetAllVehicles :many
SELECT * 
FROM vehicles 
WHERE 
  (($1 = '') OR (provider = $1)) AND archived_at IS NULL
ORDER BY water_capacity DESC
LIMIT $2 OFFSET $3;

-- name: GetAllVehiclesCount :one
SELECT COUNT(*) 
  FROM vehicles 
  WHERE 
    ($1 = '' OR provider = $1)
    AND archived_at IS NULL;

-- name: GetAllVehiclesWithArchived :many
SELECT * FROM vehicles 
WHERE 
 (($1 = '') OR (provider = $1))
ORDER BY water_capacity DESC
LIMIT $2 OFFSET $3;

-- name: GetAllVehiclesWithArchivedCount :one
SELECT COUNT(*) 
  FROM vehicles 
  WHERE 
    ($1 = '' OR provider = $1);

-- name: GetAllVehiclesByType :many
SELECT * 
FROM vehicles 
WHERE 
  type = $1
  AND ($2 = '' OR provider = $2)
  AND archived_at IS NULL
ORDER BY water_capacity DESC
LIMIT $3 OFFSET $4;

-- name: GetAllVehiclesByTypeCount :one
SELECT COUNT(*) 
  FROM vehicles 
  WHERE 
    type = $1
    AND ($2 = '' OR provider = $2)
    AND archived_at IS NULL;

-- name: GetAllVehiclesByTypeWithArchived :many
SELECT * 
FROM vehicles 
WHERE 
  type = $1
  AND ($2 = '' OR provider = $2)
ORDER BY water_capacity DESC
LIMIT $3 OFFSET $4;

-- name: GetAllVehiclesByTypeWithArchivedCount :one
SELECT COUNT(*) 
  FROM vehicles 
  WHERE 
    type = $1
    AND ($2 = '' OR provider = $2);

-- name: GetAllArchivedVehicles :many
SELECT * 
FROM vehicles 
WHERE 
  archived_at IS NOT NULL
ORDER BY water_capacity DESC;

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

-- name: ArchiveVehicle :one
UPDATE vehicles SET archived_at = $2 WHERE id = $1 RETURNING id;

-- name: DeleteVehicle :one
DELETE FROM vehicles WHERE id = $1 RETURNING id;

-- name: GetAllVehiclesWithWateringPlanCount :many
SELECT 
    v.number_plate,
    COUNT(vwp.watering_plan_id) AS watering_plan_count
FROM vehicles v
LEFT JOIN vehicle_watering_plans vwp ON v.id = vwp.vehicle_id
GROUP BY v.number_plate
ORDER BY watering_plan_count DESC;


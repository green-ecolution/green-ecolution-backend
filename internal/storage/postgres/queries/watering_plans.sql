-- name: GetAllWateringPlans :many
SELECT * FROM watering_plans;

-- name: GetWateringPlanByID :one
SELECT * FROM watering_plans WHERE id = $1;

-- name: GetTransporterByWateringPlanID :one
SELECT v.* FROM vehicles v
JOIN 
  vehicle_watering_plans vwp ON v.id = vwp.vehicle_id
WHERE 
  vwp.watering_plan_id = $1
AND v.type = 'transporter';

-- name: GetTrailerByWateringPlanID :one
SELECT v.* FROM vehicles v
JOIN 
  vehicle_watering_plans vwp ON v.id = vwp.vehicle_id
WHERE 
  vwp.watering_plan_id = $1
AND v.type = 'trailer';

-- name: CreateWateringPlan :one
INSERT INTO watering_plans (
  date, description, watering_plan_status, distance, total_water_required
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING id;

-- name: UpdateWateringPlan :exec
UPDATE watering_plans SET
  date = $2,
  description = $3,
  watering_plan_status = $4,
  distance = $5,
  total_water_required = $6
WHERE id = $1;

-- name: DeleteWateringPlan :one
DELETE FROM watering_plans WHERE id = $1 RETURNING id;

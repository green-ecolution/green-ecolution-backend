-- name: GetAllWateringPlans :many
SELECT * FROM watering_plans;

-- name: GetWateringPlanByID :one
SELECT * FROM watering_plans WHERE id = $1;

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

-- name: GetTransporterByWateringPlanID :one
SELECT v.* FROM vehicles v
JOIN vehicle_watering_plans vwp ON v.id = vwp.vehicle_id
WHERE vwp.watering_plan_id = $1
AND v.type = 'transporter';

-- name: GetTrailerByWateringPlanID :one
SELECT v.* FROM vehicles v
JOIN vehicle_watering_plans vwp ON v.id = vwp.vehicle_id
WHERE vwp.watering_plan_id = $1
AND v.type = 'trailer';

-- name: SetVehicleToWateringPlan :exec
INSERT INTO vehicle_watering_plans (vehicle_id, watering_plan_id)
VALUES ($1, $2);

-- name: DeleteAllVehiclesWateringPlan :exec
DELETE FROM vehicle_watering_plans
WHERE watering_plan_id = $1;

-- name: GetTreeClustersByWateringPlanID :many
SELECT tc.*
FROM tree_clusters tc
JOIN tree_cluster_watering_plans tcwp ON tc.id = tcwp.tree_cluster_id
WHERE tcwp.watering_plan_id = $1;

-- name: SetTreeclusterToWateringPlan :exec
INSERT INTO tree_cluster_watering_plans (tree_cluster_id, watering_plan_id)
VALUES ($1, $2);

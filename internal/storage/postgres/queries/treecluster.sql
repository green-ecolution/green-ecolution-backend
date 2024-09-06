-- name: GetAllTreeClusters :many
SELECT * FROM tree_clusters;

-- name: GetTreeClusterByID :one
SELECT * FROM tree_clusters WHERE id = $1;

-- name: GetSensorByTreeClusterID :one
SELECT sensors.* FROM sensors JOIN tree_clusters ON sensors.id = tree_clusters.sensor_id WHERE tree_clusters.id = $1;

-- name: CreateTreeCluster :one
INSERT INTO tree_clusters (
  region, address, description, moisture_level, latitude, longitude
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING id;

-- name: UpdateTreeClusterWateringStatus :exec
UPDATE tree_clusters SET
  watering_status = $2
WHERE id = $1;

-- name: UpdateTreeClusterMoistureLevel :exec
UPDATE tree_clusters SET
  moisture_level = $2
WHERE id = $1;

-- name: UpdateTreeClusterSoilCondition :exec
UPDATE tree_clusters SET
  soil_condition = $2
WHERE id = $1;

-- name: UpdateTreeClusterLastWatered :exec
UPDATE tree_clusters SET
  last_watered = $2
WHERE id = $1;

-- name: UpdateTreeClusterGeometry :exec
UPDATE tree_clusters SET
  geometry = ST_SetSRID(ST_MakePoint($2, $3), 4326)
WHERE id = $1;

-- name: UpdateTreeCluster :exec
UPDATE tree_clusters SET
  region = $2,
  address = $3,
  description = $4,
  latitude = $5,
  longitude = $6
WHERE id = $1;

-- name: ArchiveTreeCluster :exec
UPDATE tree_clusters SET
  archived = TRUE
WHERE id = $1;

-- name: DeleteTreeCluster :exec
DELETE FROM tree_clusters WHERE id = $1;

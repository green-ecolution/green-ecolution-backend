-- name: GetAllTreeClusters :many
SELECT * FROM tree_clusters;

-- name: GetTreeClusterByID :one
SELECT * FROM tree_clusters WHERE id = $1;

-- name: GetSensorByTreeClusterID :one
SELECT sensors.* FROM sensors JOIN tree_clusters ON sensors.id = tree_clusters.sensor_id WHERE tree_clusters.id = $1;

-- name: CreateTreeCluster :one
INSERT INTO tree_clusters (
  region_id, address, description, moisture_level, latitude, longitude, watering_status, soil_condition
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING id;

-- name: UpdateTreeCluster :exec
UPDATE tree_clusters SET
  region_id = $2,
  address = $3,
  description = $4,
  latitude = $5,
  longitude = $6,
  moisture_level = $7,
  watering_status = $8,
  soil_condition = $9,
  last_watered = $10,
  archived = $11
WHERE id = $1;

-- name: ArchiveTreeCluster :exec
UPDATE tree_clusters SET
  archived = TRUE
WHERE id = $1;

-- name: DeleteTreeCluster :exec
DELETE FROM tree_clusters WHERE id = $1;

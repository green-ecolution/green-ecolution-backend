-- name: GetAllTreeClusters :many
SELECT * FROM tree_clusters;

-- name: GetTreeClusterByID :one
SELECT * FROM tree_clusters WHERE id = $1;

-- name: GetSensorByTreeClusterID :one
SELECT sensors.* FROM sensors JOIN tree_clusters ON sensors.id = tree_clusters.sensor_id WHERE tree_clusters.id = $1;

-- name: GetRegionByTreeClusterID :one
SELECT regions.* FROM regions JOIN tree_clusters ON regions.id = tree_clusters.region_id WHERE tree_clusters.id = $1;

-- name: GetLinkedTreesByTreeClusterID :many
SELECT trees.* FROM trees JOIN tree_clusters ON trees.tree_cluster_id = tree_clusters.id WHERE tree_clusters.id = $1;

-- name: CreateTreeCluster :one
INSERT INTO tree_clusters (
  name, region_id, address, description, moisture_level, latitude, longitude, watering_status, soil_condition
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING id;

-- name: UpdateTreeCluster :exec
UPDATE tree_clusters SET
  name = $2,
  region_id = $3,
  address = $4,
  description = $5,
  latitude = $6,
  longitude = $7,
  moisture_level = $8,
  watering_status = $9,
  soil_condition = $10,
  last_watered = $11,
  archived = $12
WHERE id = $1;

-- name: ArchiveTreeCluster :exec
UPDATE tree_clusters SET
  archived = TRUE
WHERE id = $1;

-- name: DeleteTreeCluster :exec
DELETE FROM tree_clusters WHERE id = $1;

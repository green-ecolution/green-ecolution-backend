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
  name, region_id, address, description, moisture_level, watering_status, soil_condition
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING id;

-- name: LinkTreesToTreeCluster :exec
UPDATE trees SET tree_cluster_id = $2 WHERE id = ANY($1::int[]);

-- name: SetTreeClusterLocation :exec
UPDATE tree_clusters SET
  latitude = $2,
  longitude = $3,
  geometry = ST_SetSRID(ST_MakePoint($2, $3), 4326)
WHERE id = $1;

-- name: UpdateTreeCluster :exec
UPDATE tree_clusters SET
  name = $2,
  region_id = $3,
  address = $4,
  description = $5,
  moisture_level = $6,
  watering_status = $7,
  soil_condition = $8,
  last_watered = $9,
  archived = $10
WHERE id = $1;

-- name: ArchiveTreeCluster :exec
UPDATE tree_clusters SET
  archived = TRUE
WHERE id = $1;

-- name: DeleteTreeCluster :exec
DELETE FROM tree_clusters WHERE id = $1;


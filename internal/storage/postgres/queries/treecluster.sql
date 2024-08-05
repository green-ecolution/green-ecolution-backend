-- name: GetAllTreeClusters :many
SELECT * FROM tree_clusters;

-- name: GetTreeClusterByID :one
SELECT * FROM tree_clusters WHERE id = $1;

-- name: GetSensorByTreeClusterID :one
SELECT sensors.* FROM sensors JOIN tree_clusters ON sensors.id = tree_clusters.sensor_id WHERE tree_clusters.id = $1;

-- name: CreateTreeCluster :one
INSERT INTO tree_clusters (
  region, address, latitude, longitude, geometry
) VALUES (
  $1, $2, $3, $4, ST_GeomFromText($5, 4326)
) RETURNING id;

-- name: UpdateTreeCluster :exec
UPDATE tree_clusters SET
  region = $2,
  address = $3,
  latitude = $4,
  longitude = $5
WHERE id = $1;

-- name: ArchiveTreeCluster :exec
UPDATE tree_clusters SET
  archived = TRUE
WHERE id = $1;

-- name: DeleteTreeCluster :exec
DELETE FROM tree_clusters WHERE id = $1;

-- name: GetAllTreeClusters :many
SELECT * FROM tree_clusters 
ORDER BY name ASC
LIMIT $1 OFFSET $2;

-- name: GetAllTreeClustersCount :one
SELECT COUNT(*) FROM tree_clusters;

-- name: GetTreeClusterByID :one
SELECT * FROM tree_clusters WHERE id = $1;

-- name: GetTreesClustersByIDs :many
SELECT * FROM tree_clusters WHERE id = ANY($1::int[]);

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

-- name: RemoveTreeClusterLocation :exec
UPDATE tree_clusters SET
  latitude = NULL,
  longitude = NULL,
  geometry = NULL
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

-- name: ArchiveTreeCluster :one
UPDATE tree_clusters SET
  archived = TRUE
WHERE id = $1 RETURNING id;

-- name: DeleteTreeCluster :one
DELETE FROM tree_clusters WHERE id = $1 RETURNING id;

-- name: CalculateTreesCentroid :one
SELECT ST_AsText(ST_Centroid(ST_Collect(geometry)))::text AS centroid FROM trees WHERE trees.tree_cluster_id = $1;

-- name: GetAllLatestSensorDataByTreeClusterID :many
SELECT sd.*
FROM sensor_data sd
JOIN sensors s ON sd.sensor_id = s.id
JOIN trees t ON t.sensor_id = s.id
JOIN tree_clusters tc ON t.tree_cluster_id = tc.id
WHERE tc.id = $1
  AND sd.id = (
    SELECT id
    FROM sensor_data
    WHERE sensor_id = s.id
    ORDER BY created_at DESC
    LIMIT 1
  );

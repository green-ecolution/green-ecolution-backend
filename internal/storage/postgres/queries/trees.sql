-- name: GetAllTrees :many
SELECT * FROM trees;

-- name: GetTreeByID :one
SELECT * FROM trees WHERE id = $1;

-- name: GetTreesByIDs :many
SELECT * FROM trees WHERE id = ANY($1::int[]);

-- name: GetTreesByTreeClusterID :many
SELECT * FROM trees WHERE tree_cluster_id = $1;

-- name: GetTreeByCoordinates :one
SELECT * FROM trees WHERE latitude = $1 AND longitude = $2 LIMIT 1;

-- name: GetAllImagesByTreeID :many
SELECT images.* FROM images JOIN tree_images ON images.id = tree_images.image_id WHERE tree_images.tree_id = $1;

-- name: GetSensorByTreeID :one
SELECT sensors.* FROM sensors JOIN trees ON sensors.id = trees.sensor_id WHERE trees.id = $1;

-- name: GetTreeClusterByTreeID :one
SELECT tree_clusters.* FROM tree_clusters JOIN trees ON tree_clusters.id = trees.tree_cluster_id WHERE trees.id = $1;

-- name: CreateTree :one
INSERT INTO trees (
  tree_cluster_id, sensor_id, planting_year, species, tree_number, readonly, description, watering_status, latitude, longitude
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING id;

-- name: UpdateTree :exec
UPDATE trees SET
  tree_cluster_id = $2,
  sensor_id = $3,
  planting_year = $4,
  species = $5,
  tree_number = $6,
  readonly = $7,
  watering_status = $8,
  description = $9
WHERE id = $1;

-- name: SetTreeLocation :exec
UPDATE trees SET
  latitude = $2,
  longitude = $3,
  geometry = ST_SetSRID(ST_MakePoint($2, $3), 4326)
WHERE id = $1;

-- name: UpdateTreeClusterID :exec
UPDATE trees SET tree_cluster_id = $2 WHERE id = ANY($1::int[]);

-- name: LinkTreeImage :exec
INSERT INTO tree_images (
  tree_id, image_id
) VALUES (
  $1, $2
);

-- name: UnlinkTreeImage :exec
DELETE FROM tree_images WHERE tree_id = $1 AND image_id = $2;

-- name: UnlinkAllTreeImages :exec
DELETE FROM tree_images WHERE tree_id = $1;

-- name: UpdateTreeGeometry :exec
UPDATE trees SET
  geometry = ST_GeomFromText($2, 4326)
WHERE id = $1;

-- name: DeleteTree :exec
DELETE FROM trees WHERE id = $1;

-- name: UnlinkTreeClusterID :exec
UPDATE trees SET tree_cluster_id = NULL WHERE tree_cluster_id = $1;

-- name: CalculateGroupedCentroids :one
SELECT ST_AsText(ST_Centroid(ST_Collect(geometry)))::text AS centroid FROM trees WHERE id = ANY($1::int[]);

-- name: GetAllTrees :many
SELECT * FROM trees;

-- name: GetTreeByID :one
SELECT * FROM trees WHERE id = $1;

-- name: GetTreesByTreeClusterID :many
SELECT * FROM trees WHERE tree_cluster_id = $1;

-- name: GetAllImagesByTreeID :many
SELECT images.* FROM images JOIN tree_images ON images.id = tree_images.image_id WHERE tree_images.tree_id = $1;

-- name: GetSensorByTreeID :one
SELECT sensors.* FROM sensors JOIN trees ON sensors.id = trees.sensor_id WHERE trees.id = $1;

-- name: GetTreeClusterByTreeID :one
SELECT tree_clusters.* FROM tree_clusters JOIN trees ON tree_clusters.id = trees.tree_cluster_id WHERE trees.id = $1;

-- name: CreateTree :one
INSERT INTO trees (
  tree_cluster_id, sensor_id, age, height_above_sea_level, planting_year, species, tree_number, latitude, longitude, geometry
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, ST_GeomFromText($10, 4326)
) RETURNING id;

-- name: UpdateTree :exec
UPDATE trees SET
  tree_cluster_id = $2,
  sensor_id = $3,
  age = $4,
  height_above_sea_level = $5,
  planting_year = $6,
  species = $7,
  tree_number = $8,
  latitude = $9,
  longitude = $10,
  geometry = ST_GeomFromText($11, 4326)
WHERE id = $1;

-- name: LinkTreeImage :exec
INSERT INTO tree_images (
  tree_id, image_id
) VALUES (
  $1, $2
);

-- name: UnlinkTreeImage :exec
DELETE FROM tree_images WHERE tree_id = $1 AND image_id = $2;

-- name: UpdateTreeGeometry :exec
UPDATE trees SET
  geometry = ST_GeomFromText($2, 4326)
WHERE id = $1;

-- name: DeleteTree :exec
DELETE FROM trees WHERE id = $1;

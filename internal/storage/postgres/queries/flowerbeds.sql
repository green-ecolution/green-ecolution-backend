-- name: GetAllFlowerbeds :many
SELECT * FROM flowerbeds;

-- name: GetFlowerbedByID :one
SELECT * FROM flowerbeds WHERE id = $1;

-- name: GetSensorByFlowerbedID :one
SELECT sensors.* FROM sensors JOIN flowerbeds ON sensors.id = flowerbeds.sensor_id WHERE flowerbeds.id = $1;

-- name: GetAllImagesByFlowerbedID :many
SELECT images.* FROM images JOIN flowerbed_images ON images.id = flowerbed_images.image_id WHERE flowerbed_images.flowerbed_id = $1;

-- name: GetRegionByFlowerbedID :one
SELECT regions.* FROM regions JOIN flowerbeds ON regions.id = flowerbeds.region_id WHERE flowerbeds.id = $1;

-- name: CreateFlowerbed :one
INSERT INTO flowerbeds (
  sensor_id, size, description, number_of_plants, moisture_level, region_id, address, latitude, longitude, geometry
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, ST_GeomFromText($10, 4326)
) RETURNING id;

-- name: LinkFlowerbedImage :exec
INSERT INTO flowerbed_images (flowerbed_id, image_id) VALUES ($1, $2);

-- name: UnlinkFlowerbedImage :one
DELETE FROM flowerbed_images WHERE flowerbed_id = $1 AND image_id = $2 RETURNING flowerbed_id;

-- name: UnlinkAllFlowerbedImages :one
DELETE FROM flowerbed_images WHERE flowerbed_id = $1 RETURNING flowerbed_id;

-- name: UnlinkSensorIDFromFlowerbeds :exec
UPDATE flowerbeds SET sensor_id = NULL WHERE sensor_id = $1;

-- name: UpdateFlowerbed :exec
UPDATE flowerbeds SET
  sensor_id = $2,
  size = $3,
  description = $4,
  number_of_plants = $5,
  moisture_level = $6,
  region_id = $7,
  address = $8,
  latitude = $9,
  longitude = $10,
  geometry = ST_GeomFromText($11, 4326)
WHERE id = $1;

-- name: UpdateFlowerbedGeometry :exec
UPDATE flowerbeds SET
  geometry = ST_GeomFromText($2, 4326)
WHERE id = $1;

-- name: ArchiveFlowerbed :one
UPDATE flowerbeds SET
  archived = TRUE
WHERE id = $1 RETURNING id;

-- name: DeleteFlowerbed :one
DELETE FROM flowerbeds WHERE id = $1 RETURNING id;

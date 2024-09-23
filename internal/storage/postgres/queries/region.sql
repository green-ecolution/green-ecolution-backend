-- name: GetRegionByPoint :one
SELECT * FROM regions WHERE ST_Contains(geometry, ST_GeomFromText($1, 4326));

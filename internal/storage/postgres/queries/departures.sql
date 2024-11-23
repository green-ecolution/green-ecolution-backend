-- name: GetAllDepartures :many
SELECT * FROM departures;

-- name: GetDepartureByID :one
SELECT * FROM departures WHERE id = $1;

-- name: CreateDeparture :one
INSERT INTO departures (
  name, description, latitude, longitude
) VALUES (
  $1, $2, $3, $4
) RETURNING id;

-- name: UpdateDeparture :exec
UPDATE departures SET
  name = $2,
  description = $3,
  latitude = $4,
  longitude = $5
WHERE id = $1;

-- name: DeleteDeparture :one
DELETE FROM departures WHERE id = $1 RETURNING id;

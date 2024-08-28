-- name: GetAllImages :many
SELECT * FROM images;

-- name: GetImageByID :one
SELECT * FROM images WHERE id = $1;

-- name: CreateImage :one
INSERT INTO images (
  url, filename, mime_type
) VALUES (
  $1, $2, $3
) RETURNING id;

-- name: UpdateImage :exec
UPDATE images SET
  url = $2, filename = $3, mime_type = $4
WHERE id = $1;

-- name: DeleteImage :exec
DELETE FROM images WHERE id = $1;

-- +goose Up
-- +goose StatementBegin
INSERT INTO images (id, url) VALUES (1, 'https://avatars.githubusercontent.com/u/165842746?s=96&v=4');
INSERT INTO images (id, url, filename, mime_type) VALUES (2, 'https://app.dev.green-ecolution.de/api/v1/images/avatar.png', 'avatar.png', 'image/png');
ALTER SEQUENCE images_id_seq RESTART WITH 3;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM images WHERE url = 'https://avatars.githubusercontent.com/u/165842746?s=96&v=4';
DELETE FROM images WHERE url = 'https://app.dev.green-ecolution.de/api/v1/images/avatar.png';
ALTER SEQUENCE images_id_seq RESTART WITH 1;
-- +goose StatementEnd

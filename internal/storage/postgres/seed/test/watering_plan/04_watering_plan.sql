-- +goose Up
-- +goose StatementBegin
INSERT INTO watering_plans (id, date, description, watering_plan_status, distance, total_water_required)
VALUES 
  (1, '2024-09-22', 'New watering plan for the west side of the city', 'planned', 63.0, 6000.0),
  (2, '2024-08-03', 'New watering plan for the east side of the city', 'active', 63.0, 6000.0),
  (3, '2024-06-12', 'Very important watering plan due to no rainfall', 'finished', 63.0, 6000.0),
  (4, '2024-06-10', 'New watering plan for the south side of the city', 'not competed', 63.0, 6000.0),
  (5, '2024-06-04', 'Cancelled due to flood', 'cancelled', 63.0, 6000.0);
ALTER SEQUENCE watering_plans_id_seq RESTART WITH 6;

INSERT INTO user_watering_plans (user_id, watering_plan_id) VALUES 
('b55bd65c-301f-4e2a-9ab2-a91c6cdaed20', 1),
('c7a6434a-e91d-4d21-af07-d5c7e6b7ba5a', 1),
('2db62407-0d00-4c02-8d73-93f3d0701d49', 1);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM watering_plans;
DELETE FROM user_watering_plans;
-- +goose StatementEnd
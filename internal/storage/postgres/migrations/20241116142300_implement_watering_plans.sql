-- +goose Up
CREATE TABLE IF NOT EXISTS departures (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  name TEXT NOT NULL,
  description TEXT NOT NULL, 
  latitude FLOAT NOT NULL,
  longitude FLOAT NOT NULL
);

CREATE TYPE watering_plan_status AS ENUM ('planed', 'active', 'cancelled', 'finished', 'not competed', 'unknown');

CREATE TABLE IF NOT EXISTS watering_plans (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  date DATE NOT NULL,
  description TEXT NOT NULL,
  watering_plan_status watering_plan_status NOT NULL DEFAULT 'unknown',
  distance FLOAT,
  total_water_required FLOAT,
  departure_id INT NOT NULL,
  FOREIGN KEY (departure_id) REFERENCES departures(id)
);

CREATE TABLE IF NOT EXISTS user_watering_plans (
  user_id UUID NOT NULL,
  watering_plan_id INT NOT NULL,
  PRIMARY KEY (user_id, watering_plan_id),
  FOREIGN KEY (watering_plan_id) REFERENCES watering_plans(id)
);

CREATE TABLE IF NOT EXISTS vehicle_watering_plans (
  vehicle_id INT NOT NULL,
  watering_plan_id INT NOT NULL,
  PRIMARY KEY (vehicle_id, watering_plan_id),
  FOREIGN KEY (vehicle_id) REFERENCES vehicles(id),
  FOREIGN KEY (watering_plan_id) REFERENCES watering_plans(id)
);

CREATE TABLE IF NOT EXISTS departure_watering_plans (
  departure_id INT NOT NULL,
  watering_plan_id INT NOT NULL,
  PRIMARY KEY (departure_id, watering_plan_id),
  FOREIGN KEY (departure_id) REFERENCES departures(id),
  FOREIGN KEY (watering_plan_id) REFERENCES watering_plans(id)
);

CREATE TABLE IF NOT EXISTS tree_cluster_watering_plans (
  tree_cluster_id INT NOT NULL,
  watering_plan_id INT NOT NULL,
  PRIMARY KEY (tree_cluster_id, watering_plan_id),
  FOREIGN KEY (tree_cluster_id) REFERENCES tree_clusters(id),
  FOREIGN KEY (watering_plan_id) REFERENCES watering_plans(id)
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_updated_at_column()
  RETURNS TRIGGER
  AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$
language 'plpgsql';
-- +goose StatementEnd

CREATE TRIGGER update_departures_updated_at
BEFORE UPDATE ON departures
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_watering_plans_updated_at
BEFORE UPDATE ON watering_plans
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_departures_updated_at ON departures;
DROP TRIGGER IF EXISTS update_watering_plans_updated_at ON watering_plans;

DROP TABLE IF EXISTS user_watering_plans;
DROP TABLE IF EXISTS vehicle_watering_plans;
DROP TABLE IF EXISTS departure_watering_plans;
DROP TABLE IF EXISTS tree_cluster_watering_plans;
DROP TABLE IF EXISTS watering_plans;
DROP TABLE IF EXISTS departures;

DROP TYPE IF EXISTS watering_plan_status;
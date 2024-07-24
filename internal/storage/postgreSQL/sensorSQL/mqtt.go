package sensorSQL

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgreSQL/entities/sensor"
)

type SensorRepository struct {
	db *sql.DB
}

func NewSensorRepository(db *sql.DB) *SensorRepository {
	return &SensorRepository{db: db}
}

func (r *SensorRepository) Insert(ctx context.Context, data *sensor.MqttEntity) (*sensor.MqttEntity, error) {
	query := `
        INSERT INTO sensors (id, tree_id, data)
        VALUES ($1, $2, $3)
        ON CONFLICT (id) DO UPDATE SET
            tree_id = EXCLUDED.tree_id,
            data = EXCLUDED.data
        RETURNING id
    `
	return r.upsertSensor(ctx, query, data)
}

func (r *SensorRepository) Get(ctx context.Context, id string) (*sensor.MqttEntity, error) {
	query := `
        SELECT id, tree_id, data
        FROM sensors
        WHERE (data::jsonb -> 'end_device_ids' ->> 'device_id') = $1
    `
	return r.querySingleSensor(ctx, query, id)
}

func (r *SensorRepository) GetFirst(ctx context.Context) (*sensor.MqttEntity, error) {
	query := `
        SELECT id, tree_id, data
        FROM sensors
        ORDER BY id
        LIMIT 1
    `
	return r.querySingleSensor(ctx, query)
}

func (r *SensorRepository) GetLastByTreeID(ctx context.Context, treeID string) (*sensor.MqttEntity, error) {
	query := `
        SELECT id, tree_id, data
        FROM sensors
        WHERE tree_id = $1
        ORDER BY data->>'received_at' DESC
        LIMIT 1
    `
	return r.querySingleSensor(ctx, query, treeID)
}
func (r *SensorRepository) GetAllByTreeID(ctx context.Context, treeID string) ([]*sensor.MqttEntity, error) {
	query := `
        SELECT id, tree_id, data
        FROM sensors
        WHERE tree_id = $1
    `
	return r.queryMultipleSensors(ctx, query, treeID)
}
func (r *SensorRepository) upsertSensor(ctx context.Context, query string, data *sensor.MqttEntity) (*sensor.MqttEntity, error) {
	dataJSON, err := json.Marshal(data.Data)
	if err != nil {
		return nil, err
	}

	var id uuid.UUID
	err = r.db.QueryRowContext(ctx, query, data.ID, data.TreeID, dataJSON).Scan(&id)
	if err != nil {
		return nil, storage.ErrCannotUpsertData
	}

	data.ID = id
	return data, nil
}

func (r *SensorRepository) querySingleSensor(ctx context.Context, query string, args ...interface{}) (*sensor.MqttEntity, error) {
	var (
		entity   sensor.MqttEntity
		dataJSON []byte
	)

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&entity.ID, &entity.TreeID, &dataJSON,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrDataNotFound
		}
		return nil, err
	}

	if err := json.Unmarshal(dataJSON, &entity.Data); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *SensorRepository) queryMultipleSensors(ctx context.Context, query string, args ...interface{}) ([]*sensor.MqttEntity, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, storage.ErrDataNotFound
	}
	defer rows.Close()

	var data []*sensor.MqttEntity

	for rows.Next() {
		var (
			entity   sensor.MqttEntity
			dataJSON []byte
		)

		if err := rows.Scan(
			&entity.ID, &entity.TreeID, &dataJSON,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(dataJSON, &entity.Data); err != nil {
			return nil, err
		}

		data = append(data, &entity)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

package sensor

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

func (r *SensorRepository) GetAll(ctx context.Context) ([]*entities.Sensor, error) {
	log := logger.GetLogger(ctx)
	rows, err := r.store.GetAllSensors(ctx)
	if err != nil {
		log.Debug("failed to get sensors in db", "error", err)
		return nil, r.store.HandleError(err, sqlc.Sensor{})
	}

	data := r.mapper.FromSqlList(rows)
	for _, sn := range data {
		if err := r.store.MapSensorFields(ctx, sn); err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (r *SensorRepository) GetByID(ctx context.Context, id string) (*entities.Sensor, error) {
	log := logger.GetLogger(ctx)
	row, err := r.store.GetSensorByID(ctx, id)
	if err != nil {
		log.Debug("failed to get sensor by id in db", "error", err, "sensor_id", id)
		return nil, r.store.HandleError(err, sqlc.Sensor{})
	}

	data := r.mapper.FromSql(row)
	if err := r.store.MapSensorFields(ctx, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (r *SensorRepository) GetLatestSensorDataBySensorID(ctx context.Context, id string) (*entities.SensorData, error) {
	log := logger.GetLogger(ctx)
	data, err := r.store.GetLatestSensorDataBySensorID(ctx, id)
	if err != nil {
		log.Debug("failed to get latest sensor data by sensor id in db", "error", err, "sensor_id", id)
		return nil, r.store.HandleError(err, sqlc.Sensor{})
	}

	return data, nil
}

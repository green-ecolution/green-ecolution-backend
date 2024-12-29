package sensor

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
)

func (r *SensorRepository) GetAll(ctx context.Context) ([]*entities.Sensor, error) {
	rows, err := r.store.GetAllSensors(ctx)
	if err != nil {
		return nil, r.store.HandleError(err)
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
	row, err := r.store.GetSensorByID(ctx, id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	data := r.mapper.FromSql(row)
	if err := r.store.MapSensorFields(ctx, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (r *SensorRepository) GetLatestSensorDataBySensorID(ctx context.Context, id string) (*entities.SensorData, error) {
	data, err := r.store.GetLatestSensorDataBySensorID(ctx, id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	return data, nil
}
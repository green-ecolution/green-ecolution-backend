package sensor

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper"
	"github.com/pkg/errors"
)

func (r *SensorRepository) GetAll(ctx context.Context) ([]*entities.Sensor, error) {
	rows, err := r.store.GetAllSensors(ctx)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	data := r.mapper.FromSqlList(rows)
	for _, sn := range data {
		if err := r.mapFields(ctx, sn); err != nil {
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
	if err := r.mapFields(ctx, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (r *SensorRepository) GetLastSensorDataByID(ctx context.Context, id string) (*entities.SensorData, error) {
	row, err := r.store.GetLatestSensorDataByID(ctx, id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	domainData := r.mapper.FromSqlSensorData(row)
	data, err := mapper.MapSensorData(row.Data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to map sensor data")
	}
	domainData.Data = data

	return domainData, nil
}

func (r *SensorRepository) mapFields(ctx context.Context, sn *entities.Sensor) error {
	var err error

	sn.LatestData, err = r.GetLastSensorDataByID(ctx, sn.ID)
	if err != nil {
		return r.store.HandleError(err)
	}

	return nil
}

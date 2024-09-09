package sensor

import (
	"context"
	"encoding/json"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/pkg/errors"
)

func defaultSensor() *entities.Sensor {
	return &entities.Sensor{
		Status: entities.SensorStatusUnknown,
		Data:   make([]*entities.SensorData, 0),
	}
}

func (r *SensorRepository) Create(ctx context.Context, sFn ...entities.EntityFunc[entities.Sensor]) (*entities.Sensor, error) {
	entity := defaultSensor()
	for _, fn := range sFn {
		fn(entity)
	}

	id, err := r.createEntity(ctx, entity)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	entity.ID = id
	r.InsertSensorData(ctx, entity.Data)

	return r.GetByID(ctx, id)
}

func (r *SensorRepository) InsertSensorData(ctx context.Context, data []*entities.SensorData) ([]*entities.SensorData, error) {
	for _, d := range data {
		mqttData := r.mapper.FromDomainSensorData(d.Data)
		raw, err := json.Marshal(mqttData)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal mqtt data")
		}

		params := &sqlc.InsertSensorDataParams{
			SensorID: d.ID,
			Data:     raw,
		}

		err = r.store.InsertSensorData(ctx, params)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (r *SensorRepository) createEntity(ctx context.Context, sensor *entities.Sensor) (int32, error) {
	return r.store.CreateSensor(ctx, sqlc.SensorStatus(sensor.Status))
}

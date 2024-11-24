package sensor

import (
	"context"
	"encoding/json"

	"github.com/green-ecolution/green-ecolution-backend/internal/storage"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/pkg/errors"
)

func defaultSensor() *entities.Sensor {
	return &entities.Sensor{
		Status:    entities.SensorStatusUnknown,
		Data:      make([]*entities.SensorData, 0),
		Latitude:  0,
		Longitude: 0,
	}
}

func (r *SensorRepository) Create(ctx context.Context, sFn ...entities.EntityFunc[entities.Sensor]) (*entities.Sensor, error) {
	entity := defaultSensor()
	for _, fn := range sFn {
		fn(entity)
	}

	if err := r.validateSensorEntity(entity); err != nil {
		return nil, err
	}

	id, err := r.createEntity(ctx, entity)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	entity.ID = id
	if len(entity.Data) > 0 {
		_, err = r.InsertSensorData(ctx, entity.Data, id)
		if err != nil {
			return nil, err
		}
	}

	return r.GetByID(ctx, id)
}

func (r *SensorRepository) InsertSensorData(ctx context.Context, data []*entities.SensorData, id string) ([]*entities.SensorData, error) {
	if len(data) == 0 {
		return nil, errors.New("data cannot be empty")
	}

	if id == "" {
		return nil, r.store.HandleError(errors.New("sensor id cannot be empty"))
	}

	for _, d := range data {
		mqttData := r.mapper.FromDomainSensorData(d.Data)
		raw, err := json.Marshal(mqttData)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal mqtt data")
		}

		params := &sqlc.InsertSensorDataParams{
			SensorID: id,
			Data:     raw,
		}

		err = r.store.InsertSensorData(ctx, params)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (r *SensorRepository) createEntity(ctx context.Context, sensor *entities.Sensor) (string, error) {
	return r.store.CreateSensor(ctx, &sqlc.CreateSensorParams{
		ID:        sensor.ID,
		Status:    sqlc.SensorStatus(sensor.Status),
		Latitude:  sensor.Latitude,
		Longitude: sensor.Longitude,
	})
}

func (r *SensorRepository) validateSensorEntity(sensor *entities.Sensor) error {
	if sensor == nil {
		return errors.New("sensor is nil")
	}
	if sensor.Latitude < -90 || sensor.Latitude > 90 {
		return storage.ErrInvalidLatitude
	}
	if sensor.Longitude < -180 || sensor.Longitude > 180 {
		return storage.ErrInvalidLongitude
	}

	return nil
}

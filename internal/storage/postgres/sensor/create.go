package sensor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/green-ecolution/green-ecolution-backend/internal/storage"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/pkg/errors"
)

func defaultSensor() *entities.Sensor {
	return &entities.Sensor{
		Status:     entities.SensorStatusUnknown,
		LatestData: &entities.SensorData{},
		Latitude:   0,
		Longitude:  0,
	}
}

func (r *SensorRepository) Create(ctx context.Context, sFn ...entities.EntityFunc[entities.Sensor]) (*entities.Sensor, error) {
	entity := defaultSensor()
	for _, fn := range sFn {
		fn(entity)
	}

	sensor, _ := r.GetByID(ctx, entity.ID)
	if sensor != nil {
		return nil, errors.New("sensor with same ID already exists")
	}

	fmt.Println("hallo")

	if err := r.validateSensorEntity(entity); err != nil {
		return nil, err
	}

	id, err := r.createEntity(ctx, entity)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	fmt.Println(id)

	entity.ID = id
	if entity.LatestData != nil && entity.LatestData.Data != nil {
		err = r.InsertSensorData(ctx, entity.LatestData, id)
		if err != nil {
			return nil, err
		}
	}

	return r.GetByID(ctx, id)
}

func (r *SensorRepository) InsertSensorData(ctx context.Context, latestData *entities.SensorData, id string) error {
	if latestData == nil || latestData.Data == nil {
		return errors.New("latest data cannot be empty")
	}

	if id == "" {
		return r.store.HandleError(errors.New("sensor id cannot be empty"))
	}

	mqttData := r.mapper.FromDomainSensorData(latestData.Data)
	raw, err := json.Marshal(mqttData)
	if err != nil {
		return errors.Wrap(err, "failed to marshal mqtt data")
	}

	params := &sqlc.InsertSensorDataParams{
		SensorID: id,
		Data:     raw,
	}

	err = r.store.InsertSensorData(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func (r *SensorRepository) createEntity(ctx context.Context, sensor *entities.Sensor) (string, error) {
	id, err := r.store.CreateSensor(ctx, &sqlc.CreateSensorParams{
		ID:     sensor.ID,
		Status: sqlc.SensorStatus(sensor.Status),
	})
	if err != nil {
		return "", err
	}

	if err := r.store.SetSensorLocation(ctx, &sqlc.SetSensorLocationParams{
		ID:        id,
		Latitude:  sensor.Latitude,
		Longitude: sensor.Longitude,
	}); err != nil {
		return "", err
	}
	return id, nil
}

func (r *SensorRepository) validateSensorEntity(sensor *entities.Sensor) error {
	if sensor.ID == "" {
		return errors.New("sensor id cannot be empty")
	}
	if sensor.Latitude < -90 || sensor.Latitude > 90 || sensor.Latitude == 0 {
		return storage.ErrInvalidLatitude
	}
	if sensor.Longitude < -180 || sensor.Longitude > 180 || sensor.Longitude == 0 {
		return storage.ErrInvalidLongitude
	}

	return nil
}

package sensor

import (
	"context"

	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	storageSensor "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/sensor"
	"github.com/pkg/errors"
)

func (s *SensorService) HandleMessage(ctx context.Context, payload *domain.MqttPayload) (*domain.SensorData, error) {
	if payload == nil {
		return nil, handleError(errors.New("mqtt payload is nil"))
	}

	if err := s.validator.Struct(payload); err != nil {
		return nil, handleError(service.ErrValidation)
	}

	sensor, err := s.getSensor(ctx, payload)
	if err != nil {
		return nil, handleError(err)
	}

	data := domain.SensorData{
		Data: payload,
	}
	err = s.sensorRepo.InsertSensorData(ctx, &data, sensor.ID)
	if err != nil {
		return nil, handleError(err)
	}

	sensorData, err := s.sensorRepo.GetLatestSensorDataBySensorID(ctx, sensor.ID)
	if err != nil {
		return nil, err
	}

	return sensorData, nil
}

func (s *SensorService) getSensor(ctx context.Context, payload *domain.MqttPayload) (*domain.Sensor, error) {
	sensor, err := s.sensorRepo.GetByID(ctx, payload.Device)
	if err == nil && sensor != nil {
		if sensor.Latitude != payload.Latitude || sensor.Longitude != payload.Longitude || sensor.Status != domain.SensorStatusOnline {
			sensor, err = s.sensorRepo.Update(
				ctx,
				sensor.ID,
				storageSensor.WithLatitude(payload.Latitude),
				storageSensor.WithLongitude(payload.Longitude),
				storageSensor.WithStatus(domain.SensorStatusOnline))
			if err != nil {
				return nil, err
			}
		}
		return sensor, nil
	}
	sensor, err = s.sensorRepo.Create(ctx, storageSensor.WithSensorID(payload.Device),
		storageSensor.WithLatitude(payload.Latitude),
		storageSensor.WithLongitude(payload.Longitude),
		storageSensor.WithStatus(domain.SensorStatusOnline))
	if err != nil {
		return nil, err
	}
	return sensor, nil
}

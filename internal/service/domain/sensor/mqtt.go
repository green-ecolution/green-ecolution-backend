package sensor

import (
	"context"
	"errors"

	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	storageSensor "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/sensor"
)

func (s *SensorService) HandleMessage(ctx context.Context, payload *domain.MqttPayload) (*domain.SensorData, error) {
	log := logger.GetLogger(ctx)
	if payload == nil {
		log.Debug("mqtt payload is nil")
		return nil, errors.New("mqtt payload is nil")
	}

	if err := s.validator.Struct(payload); err != nil {
		log.Debug("failed to validate mqtt payload struct", "error", err)
		return nil, err
	}

	sensor, err := s.sensorRepo.GetByID(ctx, payload.Device)
	if err != nil {
		var entityNotFoundErr storage.ErrEntityNotFound
		if !errors.As(err, &entityNotFoundErr) {
			log.Error("failed to get sensor by id", "error", err)
			return nil, err
		}
	}

	if err != nil {
		var entityNotFoundErr storage.ErrEntityNotFound
		if !errors.As(err, &entityNotFoundErr) {
			log.Error("failed to get sensor by id", "error", err)
			return nil, err
		}
	}

	if sensor != nil {
		updatedSensor, err := s.updateSensorCoordsAndStatus(ctx, payload, sensor)
		if err != nil {
			log.Error("failed to update sensor", "error", err)
			return nil, err
		}
		sensor = updatedSensor
	} else {
		log.Info("a new sensor has joined the party! creating sensor record", "sensor_id", payload.Device, "sensor_latitude", payload.Latitude, "sensor_longitude", payload.Longitude)
		createdSensor, err := s.sensorRepo.Create(ctx, storageSensor.WithSensorID(payload.Device),
			storageSensor.WithLatitude(payload.Latitude),
			storageSensor.WithLongitude(payload.Longitude),
			storageSensor.WithStatus(domain.SensorStatusOnline),
		)
		if err != nil {
			log.Error("failed to update sensor", "error", err)
			return nil, err
		}
		sensor = createdSensor
	}

	data := domain.SensorData{
		Data: payload,
	}
	err = s.sensorRepo.InsertSensorData(ctx, &data, sensor.ID)
	if err != nil {
		log.Error("failed to insert sensor data", "sensor_id", sensor.ID, "error", err)
		return nil, err
	}

	sensorData, err := s.sensorRepo.GetLatestSensorDataBySensorID(ctx, sensor.ID)
	if err != nil {
		return nil, err
	}

	s.publishNewSensorDataEvent(ctx, sensorData)

	return sensorData, nil
}

func (s *SensorService) updateSensorCoordsAndStatus(ctx context.Context, payload *domain.MqttPayload, sensor *domain.Sensor) (*domain.Sensor, error) {
	log := logger.GetLogger(ctx)
	if sensor.Latitude != payload.Latitude || sensor.Longitude != payload.Longitude || sensor.Status != domain.SensorStatusOnline {
		updatedSensor, err := s.sensorRepo.Update(
			ctx,
			sensor.ID,
			storageSensor.WithLatitude(payload.Latitude),
			storageSensor.WithLongitude(payload.Longitude),
			storageSensor.WithStatus(domain.SensorStatusOnline))
		if err != nil {
			return nil, err
		}

		return updatedSensor, err
	}

	log.Debug("sensor don't need to update coordinates and status")
	return sensor, nil
}

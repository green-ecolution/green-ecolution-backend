package sensor

import (
	"context"

	"github.com/go-playground/validator/v10"
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	storageSensor "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/sensor"
	"github.com/pkg/errors"
)

type MqttService struct {
	sensorRepo  storage.SensorRepository
	isConnected bool
	validator   *validator.Validate
}

func NewMqttService(sensorRepository storage.SensorRepository) *MqttService {
	return &MqttService{
		sensorRepo: sensorRepository,
		validator:  validator.New()}
}

func (s *MqttService) HandleMessage(ctx context.Context, payload *domain.MqttPayload) (*domain.MqttPayload, error) {
	if payload == nil {
		return nil, handleError(errors.New("mqtt payload is nil"))
	}

	if err := s.validator.Struct(payload); err != nil {
		return nil, service.NewError(service.BadRequest, errors.Wrap(err, "validation error").Error())
	}

	sensor, err := s.sensorRepo.GetByID(ctx, payload.DeviceID)
	if err != nil {
		return nil, handleError(err)
	}

	if sensor == nil {
		return nil, handleError(storage.ErrSensorNotFound)
	}

	data := []*domain.SensorData{
		{
			Data: payload,
		},
	}
	_, err = s.sensorRepo.InsertSensorData(ctx, data, payload.DeviceID)
	if err != nil {
		return nil, handleError(err)
	}

	if sensor.Latitude != payload.Latitude || sensor.Longitude != payload.Longitude {
		_, err = s.sensorRepo.Update(
			ctx,
			sensor.ID,
			storageSensor.WithLatitude(payload.Latitude),
			storageSensor.WithLongitude(payload.Longitude))
		if err != nil {
			return nil, handleError(err)
		}
	}

	return payload, nil
}

func (s *MqttService) SetConnected(ready bool) {
	s.isConnected = ready
}

func (s *MqttService) Ready() bool {
	return s.isConnected
}

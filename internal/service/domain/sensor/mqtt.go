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

func (s *MqttService) HandleMessage(ctx context.Context, payload *domain.MqttPayload) ([]*domain.SensorData, error) {
	if payload == nil {
		return nil, handleError(errors.New("mqtt payload is nil"))
	}

	if err := s.validator.Struct(payload); err != nil {
		return nil, service.NewError(service.BadRequest, errors.Wrap(err, service.ErrValidation.Error()).Error())
	}

	sensor, err := s.getSensor(ctx, payload)
	if err != nil {
		return nil, handleError(err)
	}

	data := []*domain.SensorData{
		{
			Data: payload,
		},
	}
	_, err = s.sensorRepo.InsertSensorData(ctx, data, sensor.ID)
	if err != nil {
		return nil, handleError(err)
	}

	sensorData, err := s.sensorRepo.GetSensorDataByID(ctx, sensor.ID)
	if err != nil {
		return nil, err
	}
	return sensorData, nil
}

func (s *MqttService) getSensor(ctx context.Context, payload *domain.MqttPayload) (*domain.Sensor, error) {
	sensor, err := s.sensorRepo.GetByID(ctx, payload.DeviceID)
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
	sensor, err = s.sensorRepo.Create(ctx, storageSensor.WithSensorID(payload.DeviceID),
		storageSensor.WithLatitude(payload.Latitude),
		storageSensor.WithLongitude(payload.Longitude),
		storageSensor.WithStatus(domain.SensorStatusOnline))
	if err != nil {
		return nil, err
	}
	return sensor, nil
}

func (s *MqttService) SetConnected(ready bool) {
	s.isConnected = ready
}

func (s *MqttService) Ready() bool {
	return s.isConnected
}

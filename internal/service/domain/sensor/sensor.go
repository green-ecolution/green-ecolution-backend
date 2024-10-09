package sensor

import (
	"context"
	"errors"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
)

type SensorService struct {
	sensorRepo storage.SensorRepository
}

func NewSensorService(sensorRepo storage.SensorRepository) service.SensorService {
	return &SensorService{
		sensorRepo: sensorRepo,
	}
}

func (s *SensorService) GetAll(ctx context.Context) ([]*entities.Sensor, error) {
	sensors, err := s.sensorRepo.GetAll(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	return sensors, nil
}

func (s *SensorService) Ready() bool {
	return s.sensorRepo != nil
}

func handleError(err error) error {
	if errors.Is(err, storage.ErrEntityNotFound) {
		return service.NewError(service.NotFound, err.Error())
	}

	return service.NewError(service.InternalError, err.Error())
}

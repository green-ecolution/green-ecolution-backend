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
	treeRepo storage.TreeRepository
	flowerbedRepo storage.FlowerbedRepository
}

func NewSensorService(
	sensorRepo storage.SensorRepository,
	treeRepo storage.TreeRepository,
	flowerbedRepo storage.FlowerbedRepository,
) service.SensorService {
	return &SensorService{
		sensorRepo: sensorRepo,
		treeRepo: treeRepo,
		flowerbedRepo: flowerbedRepo,
	}
}

func (s *SensorService) GetAll(ctx context.Context) ([]*entities.Sensor, error) {
	sensors, err := s.sensorRepo.GetAll(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	return sensors, nil
}

func (s *SensorService) GetByID(ctx context.Context, id int32) (*entities.Sensor, error) {
	sensor, err := s.sensorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}

	return sensor, nil
}

func (s *SensorService) Delete(ctx context.Context, id int32) error {
	_, err := s.sensorRepo.GetByID(ctx, id)
	if err != nil {
		return handleError(err)
	}

	err = s.treeRepo.UnlinkSensorID(ctx, id)
	if err != nil {
		return handleError(err)
	}

	err = s.flowerbedRepo.UnlinkSensorID(ctx, id)
	if err != nil {
		return handleError(err)
	}

	err = s.sensorRepo.Delete(ctx, id)
	if err != nil {
		return handleError(err)
	}

	return nil
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

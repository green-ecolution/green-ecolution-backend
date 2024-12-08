package sensor

import (
	"context"

	"github.com/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/sensor"
)

type SensorService struct {
	sensorRepo    storage.SensorRepository
	treeRepo      storage.TreeRepository
	flowerbedRepo storage.FlowerbedRepository
	validator     *validator.Validate
}

func NewSensorService(
	sensorRepo storage.SensorRepository,
	treeRepo storage.TreeRepository,
	flowerbedRepo storage.FlowerbedRepository,
) service.SensorService {
	return &SensorService{
		sensorRepo:    sensorRepo,
		treeRepo:      treeRepo,
		flowerbedRepo: flowerbedRepo,
		validator:     validator.New(),
	}
}

func (s *SensorService) GetAll(ctx context.Context) ([]*entities.Sensor, error) {
	sensors, err := s.sensorRepo.GetAll(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	return sensors, nil
}

func (s *SensorService) GetByID(ctx context.Context, id string) (*entities.Sensor, error) {
	get, err := s.sensorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}

	return get, nil
}

func (s *SensorService) Create(ctx context.Context, sc *entities.SensorCreate) (*entities.Sensor, error) {
	if err := s.validator.Struct(sc); err != nil {
		return nil, service.NewError(service.BadRequest, errors.Wrap(err, "validation error").Error())
	}

	created, err := s.sensorRepo.Create(ctx,
		sensor.WithData(sc.Data),
		sensor.WithStatus(sc.Status),
	)

	if err != nil {
		return nil, handleError(err)
	}

	return created, nil
}

func (s *SensorService) Update(ctx context.Context, id string, su *entities.SensorUpdate) (*entities.Sensor, error) {
	if err := s.validator.Struct(su); err != nil {
		return nil, service.NewError(service.BadRequest, errors.Wrap(err, "validation error").Error())
	}

	_, err := s.sensorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}

	updated, err := s.sensorRepo.Update(ctx, id,
		sensor.WithData(su.Data),
		sensor.WithStatus(su.Status),
	)
	if err != nil {
		return nil, handleError(err)
	}

	return updated, nil
}

func (s *SensorService) Delete(ctx context.Context, id string) error {
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
		return service.NewError(service.NotFound, storage.ErrSensorNotFound.Error())
	}

	return service.NewError(service.InternalError, err.Error())
}

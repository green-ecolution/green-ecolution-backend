package sensor

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/tree"
	"github.com/green-ecolution/green-ecolution-backend/internal/worker"
)

type SensorService struct {
	sensorRepo    storage.SensorRepository
	treeRepo      storage.TreeRepository
	flowerbedRepo storage.FlowerbedRepository
	validator     *validator.Validate
	StatusUpdater *StatusUpdater
	eventManager  *worker.EventManager
}

func NewSensorService(
	sensorRepo storage.SensorRepository,
	treeRepo storage.TreeRepository,
	flowerbedRepo storage.FlowerbedRepository,
	eventManager *worker.EventManager,
) service.SensorService {
	return &SensorService{
		sensorRepo:    sensorRepo,
		treeRepo:      treeRepo,
		flowerbedRepo: flowerbedRepo,
		validator:     validator.New(),
		StatusUpdater: &StatusUpdater{sensorRepo: sensorRepo},
		eventManager:  eventManager,
	}
}

func (s *SensorService) publishNewSensorDataEvent(ctx context.Context, data *entities.SensorData) {
	log := logger.GetLogger(ctx)
	log.Debug("publish new event", "event", entities.EventTypeNewSensorData, "service", "SensorService")
	event := entities.NewEventSensorData(data)
	if err := s.eventManager.Publish(ctx, event); err != nil {
		log.Error("error while sending event after new sensor data received", "err", err)
	}
}

func (s *SensorService) GetAll(ctx context.Context, provider string) ([]*entities.Sensor, error) {
	log := logger.GetLogger(ctx)
	var sensors []*entities.Sensor
	var err error
	if provider != "" {
		sensors, err = s.sensorRepo.GetAllByProvider(ctx, provider)
	} else {
		sensors, err = s.sensorRepo.GetAll(ctx)
	}
	if err != nil {
		log.Debug("failed to fetch sensors", "error", err)
		return nil, service.MapError(ctx, err, service.ErrorLogEntityNotFound)
	}

	return sensors, nil
}

func (s *SensorService) GetByID(ctx context.Context, id string) (*entities.Sensor, error) {
	log := logger.GetLogger(ctx)
	get, err := s.sensorRepo.GetByID(ctx, id)
	if err != nil {
		log.Debug("failed to fetch sensor by id", "sensor_id", id, "error", err)
		return nil, service.MapError(ctx, err, service.ErrorLogEntityNotFound)
	}

	return get, nil
}

func (s *SensorService) Create(ctx context.Context, sc *entities.SensorCreate) (*entities.Sensor, error) {
	log := logger.GetLogger(ctx)
	if err := s.validator.Struct(sc); err != nil {
		log.Debug("failed to validate sensor struct to create", "error", err, "raw_sensor", fmt.Sprintf("%+v", sc))
		return nil, service.MapError(ctx, errors.Join(err, service.ErrValidation), service.ErrorLogValidation)
	}

	created, err := s.sensorRepo.Create(ctx, func(s *entities.Sensor) (bool, error) {
		s.LatestData = sc.LatestData
		s.Status = sc.Status
		return true, nil
	})

	if err != nil {
		log.Debug("failed to create sensor", "error", err, "sensor_id", sc.ID)
		return nil, service.MapError(ctx, err, service.ErrorLogAll)
	}

	log.Info("sensor created successfully", "sensor_id", created.ID)
	return created, nil
}

func (s *SensorService) Update(ctx context.Context, id string, su *entities.SensorUpdate) (*entities.Sensor, error) {
	log := logger.GetLogger(ctx)
	if err := s.validator.Struct(su); err != nil {
		log.Debug("failed to validate sensor struct to update", "error", err, "raw_sensor", fmt.Sprintf("%+v", su))
		return nil, service.MapError(ctx, errors.Join(err, service.ErrValidation), service.ErrorLogValidation)
	}

	_, err := s.sensorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, service.MapError(ctx, err, service.ErrorLogEntityNotFound)
	}

	updated, err := s.sensorRepo.Update(ctx, id, func(s *entities.Sensor) (bool, error) {
		s.LatestData = su.LatestData
		s.Status = su.Status
		return true, nil
	})

	if err != nil {
		log.Debug("failed to update sensor", "sensor_id", id, "error", err)
		return nil, service.MapError(ctx, err, service.ErrorLogAll)
	}

	log.Info("sensor updated successfully", "sensor_id", id)
	return updated, nil
}

func (s *SensorService) Delete(ctx context.Context, id string) error {
	log := logger.GetLogger(ctx)
	_, err := s.sensorRepo.GetByID(ctx, id)
	if err != nil {
		return service.MapError(ctx, err, service.ErrorLogEntityNotFound)
	}

	err = s.treeRepo.UnlinkSensorID(ctx, id)
	if err != nil {
		log.Debug("failed to unlink sensor from tree", "error", err, "sensor_id", id)
		return service.MapError(ctx, err, service.ErrorLogAll)
	}

	err = s.flowerbedRepo.UnlinkSensorID(ctx, id)
	if err != nil {
		log.Debug("failed to unlink sensor from flowerbed", "error", err, "sensor_id", id)
		return service.MapError(ctx, err, service.ErrorLogAll)
	}

	err = s.sensorRepo.Delete(ctx, id)
	if err != nil {
		log.Debug("failed to delete sensor", "error", err, "sensor_id", id)
		return service.MapError(ctx, err, service.ErrorLogAll)
	}

	return nil
}

func (s *SensorService) RunStatusUpdater(ctx context.Context, interval time.Duration) {
	s.StatusUpdater.RunStatusUpdater(ctx, interval)
}

// TODO: Not called by any function
func (s *SensorService) MapSensorToTree(ctx context.Context, sen *entities.Sensor) error {
	log := logger.GetLogger(ctx)
	if sen == nil {
		return errors.New("sensor cannot be nil")
	}

	nearestTree, err := s.treeRepo.FindNearestTree(ctx, sen.Latitude, sen.Longitude)
	if err != nil {
		log.Error("failed to calculate nearest tree", "sensor_id", sen.ID, "sensor_latitude", sen.Latitude, "sensor_longitude", sen.Longitude)
		return err
	}

	if nearestTree != nil {
		_, err = s.treeRepo.Update(ctx, nearestTree.ID, tree.WithSensor(sen))
		if err != nil {
			log.Error("failed to link sensor to nearest calculated tree", "tree_id", nearestTree.ID, "sensor_id", sen.ID, "error", err)
			return err
		}
	}

	return nil
}

func (s *SensorService) Ready() bool {
	return s.sensorRepo != nil
}

package sensor

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
)

type SensorRepository struct {
	store *store.Store
	SensorRepositoryMappers
}

type SensorRepositoryMappers struct {
	mapper mapper.InternalSensorRepoMapper
}

func NewSensorRepositoryMappers(sMapper mapper.InternalSensorRepoMapper) SensorRepositoryMappers {
	return SensorRepositoryMappers{
		mapper: sMapper,
	}
}

func NewSensorRepository(s *store.Store, mappers SensorRepositoryMappers) storage.SensorRepository {
	return &SensorRepository{
		store:                   s,
		SensorRepositoryMappers: mappers,
	}
}

func WithStatus(status entities.SensorStatus) entities.EntityFunc[entities.Sensor] {
	return func(s *entities.Sensor) {
		s.Status = status
	}
}

func WithLatestData(data *entities.SensorData) entities.EntityFunc[entities.Sensor] {
	return func(s *entities.Sensor) {
		s.LatestData = data
	}
}

func WithSensorID(sensorID string) entities.EntityFunc[entities.Sensor] {
	return func(s *entities.Sensor) {
		s.ID = sensorID
	}
}

func WithLatitude(lat float64) entities.EntityFunc[entities.Sensor] {
	return func(t *entities.Sensor) {
		t.Latitude = lat
	}
}

func WithLongitude(long float64) entities.EntityFunc[entities.Sensor] {
	return func(t *entities.Sensor) {
		t.Longitude = long
	}
}

func (r *SensorRepository) Delete(ctx context.Context, id string) error {
	log := logger.GetLogger(ctx)
	if err := r.store.DeleteSensor(ctx, id); err != nil {
		log.Error("failed to delete sensor entity in db", "error", err, "sensor_id", id)
		return err
	}

	log.Debug("sensor entity deleted successfully in db", "sensor_id", id)
	return nil
}

package sensor

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
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
	s.SetEntityType(store.Sensor)
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

func WithData(data []*entities.SensorData) entities.EntityFunc[entities.Sensor] {
	return func(s *entities.Sensor) {
		s.Data = data
	}
}

func (r *SensorRepository) Delete(ctx context.Context, id int32) error {
	return r.store.DeleteSensor(ctx, id)
}

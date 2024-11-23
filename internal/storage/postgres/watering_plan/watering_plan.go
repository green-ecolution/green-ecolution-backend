package watering_plan

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
)

type WateringPlanRepository struct {
	store *store.Store
	WateringPlanMappers
}

type WateringPlanMappers struct {
	mapper mapper.InternalWateringPlanMapper
}

func NewWateringPlanMappers(wMapper mapper.InternalWateringPlanMapper) WateringPlanMappers {
	return WateringPlanMappers{
		mapper: wMapper,
	}
}

func NewWateringPlanRepository(s *store.Store, mappers WateringPlanMappers) storage.WateringPlanRepository {
	s.SetEntityType(store.WateringPlan)
	return &WateringPlanRepository{
		store:               s,
		WateringPlanMappers: mappers,
	}
}

// Create implements storage.WateringPlanRepository.
func (w *WateringPlanRepository) Create(ctx context.Context, fn ...entities.EntityFunc[entities.WateringPlan]) (*entities.WateringPlan, error) {
	panic("unimplemented")
}

// Delete implements storage.WateringPlanRepository.
func (w *WateringPlanRepository) Delete(ctx context.Context, id int32) error {
	panic("unimplemented")
}

// GetAll implements storage.WateringPlanRepository.
func (w *WateringPlanRepository) GetAll(ctx context.Context) ([]*entities.WateringPlan, error) {
	panic("unimplemented")
}

// GetByID implements storage.WateringPlanRepository.
func (w *WateringPlanRepository) GetByID(ctx context.Context, id int32) (*entities.WateringPlan, error) {
	panic("unimplemented")
}

// Update implements storage.WateringPlanRepository.
func (w *WateringPlanRepository) Update(ctx context.Context, id int32, fn ...entities.EntityFunc[entities.WateringPlan]) (*entities.WateringPlan, error) {
	panic("unimplemented")
}

package depature

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
)

type DepartureRepository struct {
	store *store.Store
	DepartureMappers
}

type DepartureMappers struct {
	mapper mapper.InternalDepartureMapper
}

func NewDepartureMappers(rMapper mapper.InternalDepartureMapper) DepartureMappers {
	return DepartureMappers{
		mapper: rMapper,
	}
}

func NewDepartureRepository(s *store.Store, mappers DepartureMappers) storage.DepartureRepository {
	s.SetEntityType(store.Departure)
	return &DepartureRepository{
		store:            s,
		DepartureMappers: mappers,
	}
}

// GetAll implements storage.DepartureRepository.
func (d *DepartureRepository) GetAll(ctx context.Context) ([]*entities.Departure, error) {
	panic("unimplemented")
}

// GetByID implements storage.DepartureRepository.
func (d *DepartureRepository) GetByID(ctx context.Context, id int32) (*entities.Departure, error) {
	panic("unimplemented")
}

// Create implements storage.DepartureRepository.
func (d *DepartureRepository) Create(ctx context.Context, fn ...entities.EntityFunc[entities.Departure]) (*entities.Departure, error) {
	panic("unimplemented")
}

// Delete implements storage.DepartureRepository.
func (d *DepartureRepository) Delete(ctx context.Context, id int32) error {
	panic("unimplemented")
}

// Update implements storage.DepartureRepository.
func (d *DepartureRepository) Update(ctx context.Context, id int32, fn ...entities.EntityFunc[entities.Departure]) (*entities.Departure, error) {
	panic("unimplemented")
}
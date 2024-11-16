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

func WithName(name string) entities.EntityFunc[entities.Departure] {
	return func(d *entities.Departure) {
		d.Name = name
	}
}

func WithDescription(description string) entities.EntityFunc[entities.Departure] {
	return func(d *entities.Departure) {
		d.Description = description
	}
}

func WithLatitude(latitude *float64) entities.EntityFunc[entities.Departure] {
	return func(d *entities.Departure) {
		d.Latitude = latitude
	}
}

func WithLongitude(longitude *float64) entities.EntityFunc[entities.Departure] {
	return func(d *entities.Departure) {
		d.Longitude = longitude
	}
}

func (d *DepartureRepository) Delete(ctx context.Context, id int32) error {
	_, err := d.store.DeleteDeparture(ctx, id)

	if err != nil {
		return err
	}

	return nil
}

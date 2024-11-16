package depature

import (
	"context"
	"errors"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

func defaultDeparture() *entities.Departure {
	return &entities.Departure{
		Name: "",
		Description: "",
		Latitude: nil,
		Longitude: nil,
	}
}

func (d *DepartureRepository) Create(ctx context.Context, dFn ...entities.EntityFunc[entities.Departure]) (*entities.Departure, error) {
	entity := defaultDeparture()

	for _, fn := range dFn {
		fn(entity)
	}

	if err := d.validateDepartureEntity(entity); err != nil {
		return nil, err
	}

	id, err := d.createEntity(ctx, entity)
	if err != nil {
		return nil, d.store.HandleError(err)
	}

	entity.ID = *id
	return d.GetByID(ctx, *id)
}

func (d *DepartureRepository) createEntity(ctx context.Context, entity *entities.Departure) (*int32, error) {
	args := sqlc.CreateDepartureParams{
		Name: entity.Name,
		Description: entity.Description,
		Latitude: *entity.Latitude,
		Longitude: *entity.Longitude,
	}

	id, err := d.store.CreateDeparture(ctx, &args)
	if err != nil {
		return nil, d.store.HandleError(err)
	}

	return &id, nil
}

func (d *DepartureRepository) validateDepartureEntity(dp *entities.Departure) error {
	if dp == nil {
		return errors.New("departure is nil")
	}

	if dp.Longitude == nil {
		return errors.New("departure longitude is empty")
	}

	if dp.Latitude == nil {
		return errors.New("departure latitude is empty")
	}

	if dp.Name == "" {
		return errors.New("departure name is empty")
	}

	return nil
}
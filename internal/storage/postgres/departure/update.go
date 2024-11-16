package depature

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

func (d *DepartureRepository) Update(ctx context.Context, id int32, dFn ...entities.EntityFunc[entities.Departure]) (*entities.Departure, error) {
	entity, err := d.GetByID(ctx, id)
	if err != nil {
		return nil, d.store.HandleError(err)
	}

	for _, fn := range dFn {
		fn(entity)
	}

	if err := d.validateDepartureEntity(entity); err != nil {
		return nil, err
	}

	if err := d.updateEntity(ctx, entity); err != nil {
		return nil, err
	}

	return d.GetByID(ctx, entity.ID)
}

func (d *DepartureRepository) updateEntity(ctx context.Context, entity *entities.Departure) error {
	params := sqlc.UpdateDepartureParams{
		ID:   entity.ID,
		Name: entity.Name,
		Description: entity.Description,
		Latitude: *entity.Latitude,
		Longitude: *entity.Longitude,
	}

	return d.store.UpdateDeparture(ctx, &params)
}


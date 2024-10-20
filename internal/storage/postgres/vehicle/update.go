package vehicle

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

func (r *VehicleRepository) Update(ctx context.Context, id int32, vFn ...entities.EntityFunc[entities.Vehicle]) (*entities.Vehicle, error) {
	entity, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	for _, fn := range vFn {
		fn(entity)
	}

	if err := r.updateEntity(ctx, entity); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, entity.ID)
}

func (r *VehicleRepository) updateEntity(ctx context.Context, vehicle *entities.Vehicle) error {
	params := sqlc.UpdateVehicleParams{
		ID:            vehicle.ID,
		NumberPlate:   vehicle.NumberPlate,
		Description:   vehicle.Description,
		WaterCapacity: vehicle.WaterCapacity,
	}

	return r.store.UpdateVehicle(ctx, &params)
}

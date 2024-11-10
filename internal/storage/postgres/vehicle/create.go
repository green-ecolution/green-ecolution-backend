package vehicle

import (
	"context"
	"errors"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

func defaultVehicle() *entities.Vehicle {
	return &entities.Vehicle{
		NumberPlate:   "",
		Description:   "",
		WaterCapacity: 0,
		Type: entities.VehicleTypeUnknown,
		Status:  entities.VehicleStatusUnknown,
	}
}

func (r *VehicleRepository) Create(ctx context.Context, vFn ...entities.EntityFunc[entities.Vehicle]) (*entities.Vehicle, error) {
	entity := defaultVehicle()
	for _, fn := range vFn {
		fn(entity)
	}

	if entity.WaterCapacity == 0 {
		return nil, errors.New("water capacity is required and can not be 0")
	}

	if entity.NumberPlate == "" {
		return nil, errors.New("number plate is required")
	}

	id, err := r.createEntity(ctx, entity)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	entity.ID = *id
	return r.GetByID(ctx, *id)
}

func (r *VehicleRepository) createEntity(ctx context.Context, entity *entities.Vehicle) (*int32, error) {
	args := sqlc.CreateVehicleParams{
		NumberPlate:   entity.NumberPlate,
		Description:   entity.Description,
		WaterCapacity: entity.WaterCapacity,
		Type: sqlc.VehicleType(entity.Type),
		Status: sqlc.VehicleStatus(entity.Status),
	}

	id, err := r.store.CreateVehicle(ctx, &args)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	return &id, nil
}

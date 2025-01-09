package vehicle

import (
	"context"
	"errors"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	store "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
)

func defaultVehicle() *entities.Vehicle {
	return &entities.Vehicle{
		NumberPlate:    "",
		Description:    "",
		WaterCapacity:  0,
		Type:           entities.VehicleTypeUnknown,
		Status:         entities.VehicleStatusUnknown,
		Model:          "",
		DrivingLicense: entities.DrivingLicenseCar,
		Height:         0,
		Length:         0,
		Width:          0,
		Weight:         0,
	}
}

func (r *VehicleRepository) Create(ctx context.Context, createFn func(*entities.Vehicle) (bool, error)) (*entities.Vehicle, error) {
	if createFn == nil {
		return nil, errors.New("createFn is nil")
	}

	var createdVh *entities.Vehicle
	err := r.store.WithTx(ctx, func(s *store.Store) error {
		oldStore := r.store
		defer func() {
			r.store = oldStore
		}()
		r.store = s

		entity := defaultVehicle()
		created, err := createFn(entity)
		if err != nil {
			return err
		}

		if !created {
			return nil
		}

		if err := r.validateVehicle(entity); err != nil {
			return err
		}

		id, err := r.createEntity(ctx, entity)
		if err != nil {
			return err
		}
		createdVh, err = r.GetByID(ctx, *id)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdVh, nil
}

func (r *VehicleRepository) createEntity(ctx context.Context, entity *entities.Vehicle) (*int32, error) {
	args := sqlc.CreateVehicleParams{
		NumberPlate:    entity.NumberPlate,
		Description:    entity.Description,
		WaterCapacity:  entity.WaterCapacity,
		Type:           sqlc.VehicleType(entity.Type),
		Status:         sqlc.VehicleStatus(entity.Status),
		DrivingLicense: sqlc.DrivingLicense(entity.DrivingLicense),
		Model:          entity.Model,
		Width:          entity.Width,
		Height:         entity.Height,
		Length:         entity.Length,
		Weight:         entity.Weight,
	}

	id, err := r.store.CreateVehicle(ctx, &args)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	return &id, nil
}

func (r *VehicleRepository) validateVehicle(entity *entities.Vehicle) error {
	if entity.WaterCapacity == 0 {
		return errors.New("water capacity is required and can not be 0")
	}

	if entity.Length == 0 || entity.Width == 0 || entity.Height == 0 || entity.Weight == 0 {
		return errors.New("size measurements are required and can not be 0")
	}

	if entity.NumberPlate == "" {
		return errors.New("number plate is required")
	}

	if entity.DrivingLicense == "" {
		return errors.New("driving license is required and should be either B, BE or C")
	}

	return nil
}

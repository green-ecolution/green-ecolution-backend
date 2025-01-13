package vehicle

import (
	"context"
	"errors"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	store "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
)

func (r *VehicleRepository) Update(ctx context.Context, id int32, updateFn func(*entities.Vehicle) (bool, error)) error {
	return r.store.WithTx(ctx, func(s *store.Store) error {
		oldStore := r.store
		defer func() {
			r.store = oldStore
		}()
		r.store = s

		vh, err := r.GetByID(ctx, id)
		if err != nil {
			return err
		}

		if updateFn == nil {
			return errors.New("updateFn is nil")
		}

		updated, err := updateFn(vh)
		if err != nil {
			return err
		}

		if !updated {
			return nil
		}

		if err := r.validateVehicle(vh); err != nil {
			return err
		}

		return r.updateEntity(ctx, vh)
	})
}

func (r *VehicleRepository) updateEntity(ctx context.Context, vehicle *entities.Vehicle) error {
	params := sqlc.UpdateVehicleParams{
		ID:             vehicle.ID,
		NumberPlate:    vehicle.NumberPlate,
		Description:    vehicle.Description,
		WaterCapacity:  vehicle.WaterCapacity,
		Type:           sqlc.VehicleType(vehicle.Type),
		Status:         sqlc.VehicleStatus(vehicle.Status),
		DrivingLicense: sqlc.DrivingLicense(vehicle.DrivingLicense),
		Model:          vehicle.Model,
		Height:         vehicle.Height,
		Length:         vehicle.Length,
		Width:          vehicle.Width,
		Weight:         vehicle.Weight,
	}

	return r.store.UpdateVehicle(ctx, &params)
}

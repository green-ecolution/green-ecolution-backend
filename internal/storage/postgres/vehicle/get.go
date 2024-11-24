package vehicle

import (
	"context"
	"errors"
	"fmt"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/jackc/pgx/v5"
)

func (r *VehicleRepository) GetAll(ctx context.Context) ([]*entities.Vehicle, error) {
	rows, err := r.store.GetAllVehicles(ctx)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	return r.mapper.FromSqlList(rows), nil
}

func (r *VehicleRepository) GetByID(ctx context.Context, id int32) (*entities.Vehicle, error) {
	row, err := r.store.GetVehicleByID(ctx, id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	return r.mapper.FromSql(row), nil
}

func (r *VehicleRepository) GetByPlate(ctx context.Context, plate string) (*entities.Vehicle, error) {
	row, err := r.store.GetVehicleByPlate(ctx, plate)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	return r.mapper.FromSql(row), nil
}

func (r *VehicleRepository) GetByWateringPlan(ctx context.Context, wateringPlanID int32, vehicleType entities.VehicleType) (*entities.Vehicle, error) {
	_, err := r.store.GetWateringPlanByID(ctx, wateringPlanID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrWateringPlanNotFound
		} else {
			return nil, r.store.HandleError(err)
		}
	}

	var row *sqlc.Vehicle
	switch vehicleType {
		case entities.VehicleTypeTrailer:
			row, err = r.store.GetVehicleTrailerByWateringPlan(ctx, wateringPlanID)
		case entities.VehicleTypeTransporter:
			row, err = r.store.GetVehicleTransporterByWateringPlan(ctx, wateringPlanID)
		default:
			return nil, fmt.Errorf("unsupported vehicle type: %v", vehicleType)
	}

	if err != nil {
		return nil, r.store.HandleError(err)
	}

	return r.mapper.FromSql(row), nil
}

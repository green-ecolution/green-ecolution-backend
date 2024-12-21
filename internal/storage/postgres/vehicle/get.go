package vehicle

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

func (r *VehicleRepository) GetAll(ctx context.Context) ([]*entities.Vehicle, error) {
	rows, err := r.store.GetAllVehicles(ctx)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	return r.mapper.FromSqlList(rows), nil
}

func (r *VehicleRepository) GetAllByType(ctx context.Context, vehicleType entities.VehicleType) ([]*entities.Vehicle, error) {
    rows, err := r.store.GetAllVehiclesByType(ctx, sqlc.VehicleType(vehicleType))
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

package wateringplan

import (
	"context"
	"errors"
	"fmt"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

func (w *WateringPlanRepository) GetAll(ctx context.Context) ([]*entities.WateringPlan, error) {
	rows, err := w.store.GetAllWateringPlans(ctx)
	if err != nil {
		return nil, w.store.HandleError(err)
	}

	data := w.mapper.FromSqlList(rows)
	for _, wp := range data {
		wp.Transporter, err = w.GetLinkedVehicleByID(ctx, wp.ID, entities.VehicleTypeTransporter)
		if err != nil {
			return nil, w.store.HandleError(err)
		}

		wp.Trailer, err = w.GetLinkedVehicleByID(ctx, wp.ID, entities.VehicleTypeTrailer)
		if err != nil {
		 	if !errors.Is(err, storage.ErrEntityNotFound) {
		 		return nil, w.store.HandleError(err)
			}
		}
	}

	// TODO: get mapped data like users, treecluster
	return data, nil
}

func (w *WateringPlanRepository) GetByID(ctx context.Context, id int32) (*entities.WateringPlan, error) {
	row, err := w.store.GetWateringPlanByID(ctx, id)
	if err != nil {
		return nil, w.store.HandleError(err)
	}

	// TODO: get mapped data like users, vehicles, treecluster
	return w.mapper.FromSql(row), nil
}

func (w *WateringPlanRepository) GetLinkedVehicleByID(ctx context.Context, id int32, vehicleType entities.VehicleType) (*entities.Vehicle, error) {
	_, err := w.GetByID(ctx, id)
	if err != nil {
		return nil, w.store.HandleError(err)
	}

	var row *sqlc.Vehicle
	switch vehicleType {
		case entities.VehicleTypeTrailer:
			row, err = w.store.GetTrailerByWateringPlanID(ctx, id)
		case entities.VehicleTypeTransporter:
			row, err = w.store.GetTransporterByWateringPlanID(ctx, id)
		default:
			return nil, fmt.Errorf("unsupported vehicle type: %v", vehicleType)
	}

	if err != nil {
		return nil, w.store.HandleError(err)
	}

	return w.vehicleMapper.FromSql(row), nil
}

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
		if err := w.mapFields(ctx, wp); err != nil {
			return nil, err
		}
	}

	// TODO: get mapped data like users
	return data, nil
}

func (w *WateringPlanRepository) GetByID(ctx context.Context, id int32) (*entities.WateringPlan, error) {
	row, err := w.store.GetWateringPlanByID(ctx, id)
	if err != nil {
		return nil, w.store.HandleError(err)
	}

	wp := w.mapper.FromSql(row)
	if err := w.mapFields(ctx, wp); err != nil {
		return nil, err
	}

	// TODO: get mapped data like users
	return wp, nil
}

func (w *WateringPlanRepository) GetLinkedVehicleByIDAndType(ctx context.Context, id int32, vehicleType entities.VehicleType) (*entities.Vehicle, error) {
	var row *sqlc.Vehicle
	var err error

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

func (w *WateringPlanRepository) GetLinkedTreeClustersByID(ctx context.Context, id int32) ([]*entities.TreeCluster, error) {
	rows, err := w.store.GetTreeClustersByWateringPlanID(ctx, id)
	if err != nil {
		return nil, w.store.HandleError(err)
	}

	tc := w.clusterMapper.FromSqlList(rows)
	for _, cluster := range tc {
		if err := w.store.MapClusterFields(ctx, cluster); err != nil {
			return nil, w.store.HandleError(err)
		}
	}

	return tc, nil
}

func (w *WateringPlanRepository) GetEvaluationValues(ctx context.Context, id int32) ([]*entities.EvaluationValue, error) {
	rows, err := w.store.GetAllTreeClusterWateringPlanByID(ctx, id)
	if err != nil {
		return nil, w.store.HandleError(err)
	}

	return w.mapper.EvaluationFromSqlList(rows), nil
}

func (w *WateringPlanRepository) mapFields(ctx context.Context, wp *entities.WateringPlan) error {
	var err error

	wp.TreeClusters, err = w.GetLinkedTreeClustersByID(ctx, wp.ID)
	if err != nil {
		return w.store.HandleError(err)
	}

	wp.Transporter, err = w.GetLinkedVehicleByIDAndType(ctx, wp.ID, entities.VehicleTypeTransporter)
	if err != nil {
		return w.store.HandleError(err)
	}

	wp.Trailer, err = w.GetLinkedVehicleByIDAndType(ctx, wp.ID, entities.VehicleTypeTrailer)
	if err != nil {
		if !errors.Is(err, storage.ErrEntityNotFound) {
			return w.store.HandleError(err)
		}
		wp.Trailer = nil
	}

	if wp.Status == entities.WateringPlanStatusFinished {
		wp.Evaluation, err = w.GetEvaluationValues(ctx, wp.ID)
		if err != nil {
			return w.store.HandleError(err)
		}
	}

	// TODO: map correct users
	wp.Users = []*entities.User{}

	return nil
}

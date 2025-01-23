package wateringplan

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/jackc/pgx/v5"
)

func (w *WateringPlanRepository) GetAll(ctx context.Context) ([]*entities.WateringPlan, error) {
	log := logger.GetLogger(ctx)
	rows, err := w.store.GetAllWateringPlans(ctx)
	if err != nil {
		log.Debug("failed to get watering plan entities in db", "error", err)
		return nil, w.store.HandleError(err, sqlc.WateringPlan{})
	}

	data := w.mapper.FromSqlList(rows)
	for _, wp := range data {
		if err := w.mapFields(ctx, wp); err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (w *WateringPlanRepository) GetByID(ctx context.Context, id int32) (*entities.WateringPlan, error) {
	log := logger.GetLogger(ctx)
	row, err := w.store.GetWateringPlanByID(ctx, id)
	if err != nil {
		log.Debug("failed to get watering plan entity by id in db", "error", err, "watering_plan_id", id)
		return nil, w.store.HandleError(err, sqlc.WateringPlan{})
	}

	wp := w.mapper.FromSql(row)
	if err := w.mapFields(ctx, wp); err != nil {
		return nil, err
	}

	return wp, nil
}

func (w *WateringPlanRepository) GetLinkedVehicleByIDAndType(ctx context.Context, id int32, vehicleType entities.VehicleType) (*entities.Vehicle, error) {
	log := logger.GetLogger(ctx)
	row, err := w.store.GetVehicleByWateringPlanID(ctx, &sqlc.GetVehicleByWateringPlanIDParams{
		WateringPlanID: id,
		Type:           sqlc.VehicleType(vehicleType),
	})

	if err != nil {
		log.Debug("failed to get linked vehicle entity by id and vehicle type", "error", err, "watering_plan_id", id, "vehicle_type", vehicleType)
		return nil, err
	}

	return w.vehicleMapper.FromSql(row), nil
}

func (w *WateringPlanRepository) GetLinkedTreeClustersByID(ctx context.Context, id int32) ([]*entities.TreeCluster, error) {
	log := logger.GetLogger(ctx)
	rows, err := w.store.GetTreeClustersByWateringPlanID(ctx, id)
	if err != nil {
		log.Debug("failed to get linked tree cluster entities by watering plan id", "error", err, "watering_plan_id", id)
		return nil, err
	}

	tc := w.clusterMapper.FromSqlList(rows)
	for _, cluster := range tc {
		if err := w.store.MapClusterFields(ctx, cluster); err != nil {
			return nil, err
		}
	}

	return tc, nil
}

func (w *WateringPlanRepository) GetEvaluationValues(ctx context.Context, id int32) ([]*entities.EvaluationValue, error) {
	log := logger.GetLogger(ctx)
	rows, err := w.store.GetAllTreeClusterWateringPlanByID(ctx, id)
	if err != nil {
		log.Debug("failed to get evaluation value entities", "error", err, "watering_plan_id", id)
		return nil, err
	}

	return w.mapper.EvaluationFromSqlList(rows), nil
}

func (w *WateringPlanRepository) GetLinkedUsersByID(ctx context.Context, id int32) ([]*uuid.UUID, error) {
	log := logger.GetLogger(ctx)
	pgUUIDS, err := w.store.GetUsersByWateringPlanID(ctx, id)
	if err != nil {
		log.Error("failed to get linked user entities by watering plan id", "error", err, "watering_plan_id", id)
		return nil, err
	}

	// Convert pgtype.UUID to uuid.UUID
	var userUUIDs []*uuid.UUID
	for _, pgUUID := range pgUUIDS {
		if pgUUID.Valid {
			uuidVal := uuid.UUID(pgUUID.Bytes)
			userUUIDs = append(userUUIDs, &uuidVal)
		}
	}

	return userUUIDs, nil
}

func (w *WateringPlanRepository) mapFields(ctx context.Context, wp *entities.WateringPlan) error {
	log := logger.GetLogger(ctx)
	var err error

	wp.TreeClusters, err = w.GetLinkedTreeClustersByID(ctx, wp.ID)
	if err != nil {
		log.Debug("failed to get linked tree cluster by watering plan id", "error", err, "watering_plan_id", wp.ID)
		return w.store.HandleError(err, sqlc.WateringPlan{})
	}

	wp.Transporter, err = w.GetLinkedVehicleByIDAndType(ctx, wp.ID, entities.VehicleTypeTransporter)
	if err != nil {
		log.Debug("failed to get linked transporter by watering plan id", "error", err, "watering_plan_id", wp.ID)
		return w.store.HandleError(err, sqlc.WateringPlan{})
	}

	wp.Trailer, err = w.GetLinkedVehicleByIDAndType(ctx, wp.ID, entities.VehicleTypeTrailer)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			log.Debug("failed to get linked trailer by watering plan id", "error", err, "watering_plan_id", wp.ID)
			return w.store.HandleError(err, sqlc.WateringPlan{})
		}
		wp.Trailer = nil
	}

	wp.UserIDs, err = w.GetLinkedUsersByID(ctx, wp.ID)
	if err != nil {
		log.Debug("failed to get linked users by watering plan id", "error", err, "watering_plan_id", wp.ID)
		return w.store.HandleError(err, sqlc.WateringPlan{})
	}

	// Only load evaluation values if the watering plan is set to »finished«
	if wp.Status == entities.WateringPlanStatusFinished {
		wp.Evaluation, err = w.GetEvaluationValues(ctx, wp.ID)
		if err != nil {
			log.Debug("failed to get evaluation values by watering plan id", "error", err, "watering_plan_id", wp.ID)
			return w.store.HandleError(err, sqlc.WateringPlan{})
		}
	} else {
		wp.Evaluation = []*entities.EvaluationValue{}
	}

	return nil
}

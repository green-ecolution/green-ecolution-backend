package wateringplan

import (
	"context"
	"errors"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
)

func (w *WateringPlanRepository) Update(ctx context.Context, id int32, updateFn func(*entities.WateringPlan) (bool, error)) error {
	return w.store.WithTx(ctx, func(s *store.Store) error {
		oldStore := w.store
		defer func() {
			w.store = oldStore
		}()
		w.store = s

		entity, err := w.GetByID(ctx, id)
		if err != nil {
			return w.store.HandleError(err)
		}

		if updateFn == nil {
			return errors.New("updateFn is nil")
		}

		updated, err := updateFn(entity)
		if err != nil {
			return err
		}

		if !updated {
			return nil
		}

		if err := w.validateWateringPlan(entity); err != nil {
			return err
		}

		return w.updateEntity(ctx, entity)
	})
}

func (w *WateringPlanRepository) updateEntity(ctx context.Context, entity *entities.WateringPlan) error {
	date, err := utils.TimeToPgDate(entity.Date)
	if err != nil {
		return errors.New("failed to convert date")
	}

	totalWaterRequired, err := w.calculateRequiredWater(ctx, entity.TreeClusters)
	if err != nil {
		return err
	}

	if entity.CancellationNote != "" && entity.Status != entities.WateringPlanStatusCanceled {
		return errors.New("cancellation note should be empty, as the current watering plan is not canceled")
	}

	params := sqlc.UpdateWateringPlanParams{
		ID:                 entity.ID,
		Date:               date,
		Description:        entity.Description,
		Distance:           entity.Distance,
		TotalWaterRequired: &totalWaterRequired,
		Status:             sqlc.WateringPlanStatus(entity.Status),
		CancellationNote:   entity.CancellationNote,
	}

	if err := w.store.DeleteAllVehiclesFromWateringPlan(ctx, entity.ID); err != nil {
		return w.store.HandleError(err)
	}

	if err := w.setLinkedVehicles(ctx, entity, entity.ID); err != nil {
		return w.store.HandleError(err)
	}

	if err := w.store.DeleteAllTreeClusterFromWateringPlan(ctx, entity.ID); err != nil {
		return w.store.HandleError(err)
	}

	if err := w.setLinkedTreeClusters(ctx, entity, entity.ID); err != nil {
		return w.store.HandleError(err)
	}

	if err := w.updateConsumedWaterValues(ctx, entity); err != nil {
		return w.store.HandleError(err)
	}

	// TODO: update linked users

	return w.store.UpdateWateringPlan(ctx, &params)
}

// This function updates the consumed water values for each tree cluster in a finished watering plan.
// To save the consumed water values, the watering plan must be »finished«
func (w *WateringPlanRepository) updateConsumedWaterValues(ctx context.Context, entity *entities.WateringPlan) error {
	if entity.Status != entities.WateringPlanStatusFinished || len(entity.Evaluation) == 0 {
		return nil
	}

	for _, value := range entity.Evaluation {
		if err := w.store.UpdateTreeClusterWateringPlan(ctx, &sqlc.UpdateTreeClusterWateringPlanParams{
			WateringPlanID: entity.ID,
			TreeClusterID:  value.TreeClusterID,
			ConsumedWater:  *value.ConsumedWater,
		}); err != nil {
			return w.store.HandleError(err)
		}
	}

	return nil
}

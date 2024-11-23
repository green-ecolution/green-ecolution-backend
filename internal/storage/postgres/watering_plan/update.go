package watering_plan

import (
	"context"
	"errors"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
)

func (w *WateringPlanRepository) Update(ctx context.Context, id int32, wpFn ...entities.EntityFunc[entities.WateringPlan]) (*entities.WateringPlan, error) {
	entity, err := w.GetByID(ctx, id)
	if err != nil {
		return nil, w.store.HandleError(err)
	}

	for _, fn := range wpFn {
		fn(entity)
	}

	if err := w.validateWateringPlan(entity); err != nil {
		return nil, err
	}

	if err := w.updateEntity(ctx, entity); err != nil {
		return nil, err
	}

	return w.GetByID(ctx, entity.ID)
}

func (w *WateringPlanRepository) updateEntity(ctx context.Context, entity *entities.WateringPlan) error {
	date, err := utils.TimeToPgDate(entity.Date);
	if err != nil {
		return errors.New("failed to convert date")
	}
	
	params := sqlc.UpdateWateringPlanParams{
		ID:            entity.ID,
		Date:   date,
		Description:   entity.Description,
		Distance: entity.Distance,
		TotalWaterRequired: entity.TotalWaterRequired,
		WateringPlanStatus: sqlc.WateringPlanStatus(entities.WateringPlanStatusPlanned),
	}

	// TODO: update linked vehicles, treecluster, users

	return w.store.UpdateWateringPlan(ctx, &params)
}

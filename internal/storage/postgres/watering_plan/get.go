package watering_plan

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
)

func (w *WateringPlanRepository) GetAll(ctx context.Context) ([]*entities.WateringPlan, error) {
	rows, err := w.store.GetAllWateringPlans(ctx)
	if err != nil {
		return nil, w.store.HandleError(err)
	}

	// TODO: get mapped data like users and tree
	return w.mapper.FromSqlList(rows), nil
}

func (w *WateringPlanRepository) GetByID(ctx context.Context, id int32) (*entities.WateringPlan, error) {
	row, err := w.store.GetWateringPlanByID(ctx, id)
	if err != nil {
		return nil, w.store.HandleError(err)
	}

	// TODO: get mapped data like users and tree
	return w.mapper.FromSql(row), nil
}

package wateringplan

import (
	"context"
	"errors"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
)

func defaultWateringPlan() *entities.WateringPlan {
	return &entities.WateringPlan{
		Date:               time.Time{},
		Description:        "",
		Distance:           utils.P(0.0),
		TotalWaterRequired: utils.P(0.0),
		Users:              make([]*entities.User, 0),
		Treecluster:        make([]*entities.TreeCluster, 0),
		Transporter:        nil,
		Trailer:            nil,
	}
}

func (w *WateringPlanRepository) Create(ctx context.Context, createFn func(*entities.WateringPlan) (bool, error)) (*entities.WateringPlan, error) {
	if createFn == nil {
		return nil, errors.New("createFn is nil")
	}

	var createdWp *entities.WateringPlan
	err := w.store.WithTx(ctx, func(s *store.Store) error {
		oldStore := w.store
		defer func() {
			w.store = oldStore
		}()
		w.store = s

		entity := defaultWateringPlan()
		created, err := createFn(entity)
		if err != nil {
			return err
		}

		if !created {
			return nil
		}

		if err := w.validateWateringPlan(entity); err != nil {
			return err
		}

		id, err := w.createEntity(ctx, entity)
		if err != nil {
			return err
		}

		createdWp, err = w.GetByID(ctx, *id)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdWp, nil
}

func (w *WateringPlanRepository) createEntity(ctx context.Context, entity *entities.WateringPlan) (*int32, error) {
	date, err := utils.TimeToPgDate(entity.Date)
	if err != nil {
		return nil, errors.New("failed to convert date")
	}

	args := sqlc.CreateWateringPlanParams{
		Date:               date,
		Description:        entity.Description,
		Distance:           entity.Distance,
		TotalWaterRequired: entity.TotalWaterRequired,
		WateringPlanStatus: sqlc.WateringPlanStatus(entities.WateringPlanStatusPlanned),
	}

	// TODO: Link vehicle, treecluster and users in pivot table

	id, err := w.store.CreateWateringPlan(ctx, &args)
	if err != nil {
		return nil, w.store.HandleError(err)
	}

	return &id, nil
}

func (w *WateringPlanRepository) validateWateringPlan(entity *entities.WateringPlan) error {
	if entity.Transporter == nil || entity.Transporter.Type != entities.VehicleTypeTransporter {
		return errors.New("watering plan requires a valid transporter")
	}

	if entity.Trailer != nil && entity.Trailer.Type != entities.VehicleTypeTrailer {
		return errors.New("trailer vehicle requires a vehicle of type trailer")
	}

	if len(entity.Users) == 0 {
		return errors.New("watering plan requires employees")
	}

	if len(entity.Treecluster) == 0 {
		return errors.New("watering plan requires tree cluster")
	}

	return nil
}

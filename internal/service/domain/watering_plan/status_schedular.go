package wateringplan

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
)

type StatusSchedular struct {
	wateringPlanRepo storage.WateringPlanRepository
}

func NewStatusSchedular(wateringPlanRepo storage.WateringPlanRepository) *StatusSchedular {
	return &StatusSchedular{
		wateringPlanRepo: wateringPlanRepo,
	}
}

func (s *StatusSchedular) RunStatusSchedular(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := s.updatePlannedWateringPlanStates(ctx)
			if err != nil {
				slog.Error("Failure to update watering plan status", "error", err.Error())
			}
		case <-ctx.Done():
			slog.Info("Stopping watering plan status schedular")
			return
		}
	}
}

// This function updates the status of all watering plans that are marked as »planned« and whose date has passed. 
// If the date is in the past, the status is updated to "not competed«.
func (s *StatusSchedular) updatePlannedWateringPlanStates(ctx context.Context) error {
	wateringPlans, err := s.wateringPlanRepo.GetAllByStatus(ctx, entities.WateringPlanStatusPlanned)
	if err != nil {
		return err
	}

	for _, wateringPlan := range wateringPlans {
		if wateringPlan.Date.Before(time.Now()) {
			err = s.wateringPlanRepo.Update(ctx, wateringPlan.ID, func(wp *entities.WateringPlan) (bool, error) {
				wp.Status = entities.WateringPlanStatusNotCompeted
				return true, nil
			})
			if err != nil {
				slog.Error(
					"Failed to update watering plan status to not completed", 
					"id", fmt.Sprintf("%d", wateringPlan.ID),
					"error", err.Error(),
				)
			} else {
				slog.Info("Watering plan marked as not competed", "id", wateringPlan.ID)
			}
		}
	}

	return nil
}

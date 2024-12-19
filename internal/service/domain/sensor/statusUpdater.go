package sensor

import (
	"context"
	"log/slog"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/sensor"
)

type StatusUpdater struct {
	sensorRepo storage.SensorRepository
}

func NewStatusUpdater(sensorRepo storage.SensorRepository) *StatusUpdater {
	return &StatusUpdater{
		sensorRepo: sensorRepo,
	}
}

func (s *StatusUpdater) RunStatusUpdater(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := s.updateStaleSensorStatuses(ctx)
			if err != nil {
				slog.Error("Failure to update sensor status", "error", err.Error())
			}
		case <-ctx.Done():
			slog.Info("Stopping status updater")
			return
		}
	}
}

func (s *StatusUpdater) updateStaleSensorStatuses(ctx context.Context) error {
	sensors, err := s.sensorRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	cutoffTime := time.Now().Add(-72 * time.Hour) // 3 days ago
	for _, sens := range sensors {
		if sens.UpdatedAt.Before(cutoffTime) {
			_, err = s.sensorRepo.Update(ctx, sens.ID, sensor.WithStatus(entities.SensorStatusOffline))
			if err != nil {
				slog.Error("Failed to update sensor %s to offline: %v", sens.ID, err.Error())
			} else {
				slog.Info("Sensor marked as offline due to inactivity", "id", sens.ID)
			}
		}
	}

	return nil
}

package sensor

import (
	"context"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
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
	log := logger.GetLogger(ctx)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := s.updateStaleSensorStatuses(ctx)
			if err != nil {
				log.Error("failure to update sensor status", "error", err.Error())
			}
		case <-ctx.Done():
			log.Info("stopping status updater")
			return
		}
	}
}

func (s *StatusUpdater) updateStaleSensorStatuses(ctx context.Context) error {
	log := logger.GetLogger(ctx)
	sensors, err := s.sensorRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	cutoffTime := time.Now().Add(-72 * time.Hour) // 3 days ago
	for _, sens := range sensors {
		if sens.UpdatedAt.Before(cutoffTime) {
			_, err = s.sensorRepo.Update(ctx, sens.ID, sensor.WithStatus(entities.SensorStatusOffline))
			if err != nil {
				log.Error("failed to update sensor status to offline", "sensor_id", sens.ID, "error", err, "prev_sensor_status", sens.Status)
			} else {
				log.Info("sensor marked as offline due to inactivity", "sensor_id", sens.ID, "prev_sensor_status", sens.Status)
			}
		}
	}

	return nil
}

package sensor

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	storageMock "github.com/green-ecolution/green-ecolution-backend/internal/storage/_mock"
	"github.com/stretchr/testify/mock"
)

func TestSensorService_RunStatusSchedular(t *testing.T) {
	t.Run("should update stale sensor states periodically", func(t *testing.T) {
		ctx := context.Background()
		sensorRepo := storageMock.NewMockSensorRepository(t)
		svc := NewStatusSchedular(sensorRepo)

		staleSensor := &entities.Sensor{
			ID:        "sensor-1",
			UpdatedAt: time.Now().Add(-73 * time.Hour), // 73 hours ago
		}
		recentSensor := &entities.Sensor{
			ID:        "sensor-2",
			UpdatedAt: time.Now().Add(-1 * time.Hour), // 1 hour ago
		}

		sensorRepo.EXPECT().GetAll(mock.Anything).Return([]*entities.Sensor{staleSensor, recentSensor}, nil)
		sensorRepo.EXPECT().Update(mock.Anything, staleSensor.ID, mock.Anything).Return(staleSensor, nil)

		go func() {
			svc.RunStatusSchedular(ctx, 10*time.Millisecond)
		}()

		time.Sleep(100 * time.Millisecond)

		sensorRepo.AssertCalled(t, "GetAll", mock.Anything)
		sensorRepo.AssertCalled(t, "Update", mock.Anything, staleSensor.ID, mock.Anything)
		sensorRepo.AssertExpectations(t) // Verifies all expectations are met
	})

	t.Run("should stop updating when context is canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		sensorRepo := storageMock.NewMockSensorRepository(t)
		svc := NewStatusSchedular(sensorRepo)

		// No GetAll or Update expected, since context will be canceled early
		go func() {
			svc.RunStatusSchedular(ctx, 10*time.Millisecond)
		}()

		time.Sleep(20 * time.Millisecond)

		sensorRepo.AssertNotCalled(t, "GetAll")
		sensorRepo.AssertNotCalled(t, "Update")
		sensorRepo.AssertExpectations(t)
	})

	t.Run("should handle error from GetAll", func(t *testing.T) {
		// given
		ctx := context.Background()
		sensorRepo := storageMock.NewMockSensorRepository(t)
		svc := NewStatusSchedular(sensorRepo)

		sensorRepo.EXPECT().GetAll(mock.Anything).Return(nil, errors.New("db error"))

		go func() {
			svc.RunStatusSchedular(ctx, 10*time.Millisecond) // Run every 10ms
		}()

		time.Sleep(50 * time.Millisecond)

		sensorRepo.AssertCalled(t, "GetAll", mock.Anything)
		sensorRepo.AssertNotCalled(t, "Update")
		sensorRepo.AssertExpectations(t)
	})

	t.Run("should log error when sensor update fails", func(t *testing.T) {
		// given
		ctx := context.Background()
		sensorRepo := storageMock.NewMockSensorRepository(t)
		svc := NewStatusSchedular(sensorRepo)

		staleSensor := &entities.Sensor{
			ID:        "sensor-1",
			UpdatedAt: time.Now().Add(-73 * time.Hour), // 73 hours ago
		}

		sensorRepo.EXPECT().GetAll(mock.Anything).Return([]*entities.Sensor{staleSensor}, nil)
		sensorRepo.EXPECT().Update(mock.Anything, staleSensor.ID, mock.Anything).Return(nil, errors.New("update failed"))

		go func() {
			svc.RunStatusSchedular(ctx, 10*time.Millisecond)
		}()

		time.Sleep(50 * time.Millisecond)

		sensorRepo.AssertCalled(t, "GetAll", mock.Anything)
		sensorRepo.AssertCalled(t, "Update", mock.Anything, staleSensor.ID, mock.Anything)
		sensorRepo.AssertExpectations(t)
	})
}

package wateringplan

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	storageMock "github.com/green-ecolution/green-ecolution-backend/internal/storage/_mock"
	"github.com/stretchr/testify/mock"
)

func TestWateringPlanService_RunStatusSchedular(t *testing.T) {
	t.Run("should check for not competed watering plans periodically", func(t *testing.T) {
		ctx := context.Background()
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		svc := NewStatusSchedular(wateringPlanRepo)

		staleWateringPlan := &entities.WateringPlan{
			ID:     1,
			Date:   time.Now().Add(-73 * time.Hour), // 73 hours ago
			Status: entities.WateringPlanStatusPlanned,
		}
		recentWateringPlan := &entities.WateringPlan{
			ID:     1,
			Date:   time.Now().Add(73 * time.Hour), // 73 hours ahead
			Status: entities.WateringPlanStatusPlanned,
		}

		wateringPlanRepo.EXPECT().GetAllByStatus(
			mock.Anything,
			entities.WateringPlanStatusPlanned,
		).Return([]*entities.WateringPlan{staleWateringPlan, recentWateringPlan}, nil)

		wateringPlanRepo.EXPECT().Update(mock.Anything, staleWateringPlan.ID, mock.Anything).Return(nil)

		go func() {
			svc.RunStatusSchedular(ctx, 10*time.Millisecond)
		}()

		time.Sleep(100 * time.Millisecond)

		wateringPlanRepo.AssertCalled(t, "GetAllByStatus", mock.Anything, entities.WateringPlanStatusPlanned)
		wateringPlanRepo.AssertCalled(t, "Update", mock.Anything, staleWateringPlan.ID, mock.Anything)
		wateringPlanRepo.AssertExpectations(t)
	})

	t.Run("should stop updating when context is canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		svc := NewStatusSchedular(wateringPlanRepo)

		go func() {
			svc.RunStatusSchedular(ctx, 10*time.Millisecond)
		}()

		time.Sleep(20 * time.Millisecond)

		wateringPlanRepo.AssertNotCalled(t, "GetAllByStatus")
		wateringPlanRepo.AssertNotCalled(t, "Update")
		wateringPlanRepo.AssertExpectations(t)
	})

	t.Run("should handle error from GetAll", func(t *testing.T) {
		// given
		ctx := context.Background()
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		svc := NewStatusSchedular(wateringPlanRepo)

		wateringPlanRepo.EXPECT().GetAllByStatus(
			mock.Anything,
			entities.WateringPlanStatusPlanned,
		).Return(nil, errors.New("db error"))

		go func() {
			svc.RunStatusSchedular(ctx, 10*time.Millisecond)
		}()

		time.Sleep(50 * time.Millisecond)

		wateringPlanRepo.AssertCalled(t, "GetAllByStatus", mock.Anything, entities.WateringPlanStatusPlanned)
		wateringPlanRepo.AssertNotCalled(t, "Update")
		wateringPlanRepo.AssertExpectations(t)
	})

	t.Run("should log error when sensor update fails", func(t *testing.T) {
		// given
		ctx := context.Background()
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		svc := NewStatusSchedular(wateringPlanRepo)

		staleWateringPlan := &entities.WateringPlan{
			ID:     1,
			Date:   time.Now().Add(-73 * time.Hour), // 73 hours ago
			Status: entities.WateringPlanStatusPlanned,
		}

		wateringPlanRepo.EXPECT().GetAllByStatus(
			mock.Anything,
			entities.WateringPlanStatusPlanned,
		).Return([]*entities.WateringPlan{staleWateringPlan}, nil)
		wateringPlanRepo.EXPECT().Update(mock.Anything, staleWateringPlan.ID, mock.Anything).Return(errors.New("update failed"))

		go func() {
			svc.RunStatusSchedular(ctx, 10*time.Millisecond)
		}()

		time.Sleep(50 * time.Millisecond)

		wateringPlanRepo.AssertCalled(t, "GetAllByStatus", mock.Anything, entities.WateringPlanStatusPlanned)
		wateringPlanRepo.AssertCalled(t, "Update", mock.Anything, staleWateringPlan.ID, mock.Anything)
		wateringPlanRepo.AssertExpectations(t)
	})
}

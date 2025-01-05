package treecluster

import (
	"context"
	"testing"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	storageMock "github.com/green-ecolution/green-ecolution-backend/internal/storage/_mock"
	"github.com/green-ecolution/green-ecolution-backend/internal/worker"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestTreeClusterService_HandleUpdateWateringPlan(t *testing.T) {
	t.Run("should update tree cluster last watered", func(t *testing.T) {
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		eventManager := worker.NewEventManager(entities.EventTypeUpdateTreeCluster)
		svc := NewTreeClusterService(clusterRepo, treeRepo, regionRepo, eventManager)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go eventManager.Run(ctx)

		date := time.Date(2024, 11, 22, 0, 0, 0, 0, time.UTC)
		prevTc := entities.TreeCluster{
			ID:          1,
			LastWatered: nil,
		}
		prevWp := entities.WateringPlan{
			ID:           1,
			TreeClusters: []*entities.TreeCluster{&prevTc},
			Status:       entities.WateringPlanStatusActive,
			Date:         date,
		}

		updatedWp := entities.WateringPlan{
			ID:           1,
			TreeClusters: []*entities.TreeCluster{&prevTc},
			Status:       entities.WateringPlanStatusFinished,
			Date:         date,
		}

		event := entities.NewEventUpdateWateringPlan(&prevWp, &updatedWp)

		clusterRepo.EXPECT().Update(mock.Anything, int32(1), mock.Anything).Return(nil)

		// when
		err := svc.HandleUpdateWateringPlan(context.Background(), &event)

		// then
		assert.NoError(t, err)
	})

	t.Run("should not update tree cluster if status has not changed", func(t *testing.T) {
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		eventManager := worker.NewEventManager(entities.EventTypeUpdateTreeCluster)
		svc := NewTreeClusterService(clusterRepo, treeRepo, regionRepo, eventManager)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go eventManager.Run(ctx)

		date := time.Date(2024, 11, 22, 0, 0, 0, 0, time.UTC)
		prevTc := entities.TreeCluster{
			ID:          1,
			LastWatered: nil,
		}
		prevWp := entities.WateringPlan{
			ID:           1,
			TreeClusters: []*entities.TreeCluster{&prevTc},
			Status:       entities.WateringPlanStatusActive,
			Date:         date,
		}

		updatedWp := entities.WateringPlan{
			ID:           1,
			TreeClusters: []*entities.TreeCluster{&prevTc},
			Status:       entities.WateringPlanStatusActive,
			Date:         date,
		}

		event := entities.NewEventUpdateWateringPlan(&prevWp, &updatedWp)

		// when
		err := svc.HandleUpdateWateringPlan(context.Background(), &event)

		// then
		assert.NoError(t, err)
		clusterRepo.AssertNotCalled(t, "Update")
	})

	t.Run("should not update tree cluster if new status is not »finished«", func(t *testing.T) {
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		eventManager := worker.NewEventManager(entities.EventTypeUpdateTreeCluster)
		svc := NewTreeClusterService(clusterRepo, treeRepo, regionRepo, eventManager)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go eventManager.Run(ctx)

		date := time.Date(2024, 11, 22, 0, 0, 0, 0, time.UTC)
		prevTc := entities.TreeCluster{
			ID:          1,
			LastWatered: nil,
		}
		prevWp := entities.WateringPlan{
			ID:           1,
			TreeClusters: []*entities.TreeCluster{&prevTc},
			Status:       entities.WateringPlanStatusActive,
			Date:         date,
		}

		updatedWp := entities.WateringPlan{
			ID:           1,
			TreeClusters: []*entities.TreeCluster{&prevTc},
			Status:       entities.WateringPlanStatusCanceled,
			Date:         date,
		}

		event := entities.NewEventUpdateWateringPlan(&prevWp, &updatedWp)

		// when
		err := svc.HandleUpdateWateringPlan(context.Background(), &event)

		// then
		assert.NoError(t, err)
		clusterRepo.AssertNotCalled(t, "Update")
	})

	t.Run("should not update tree cluster if date is not the same", func(t *testing.T) {
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		eventManager := worker.NewEventManager(entities.EventTypeUpdateTreeCluster)
		svc := NewTreeClusterService(clusterRepo, treeRepo, regionRepo, eventManager)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go eventManager.Run(ctx)

		date := time.Date(2024, 11, 22, 0, 0, 0, 0, time.UTC)
		prevTc := entities.TreeCluster{
			ID:          1,
			LastWatered: nil,
		}
		prevWp := entities.WateringPlan{
			ID:           1,
			TreeClusters: []*entities.TreeCluster{&prevTc},
			Status:       entities.WateringPlanStatusActive,
			Date:         date,
		}

		updatedWp := entities.WateringPlan{
			ID:           1,
			TreeClusters: []*entities.TreeCluster{&prevTc},
			Status:       entities.WateringPlanStatusCanceled,
			Date:         time.Date(2025, 11, 22, 0, 0, 0, 0, time.UTC),
		}

		event := entities.NewEventUpdateWateringPlan(&prevWp, &updatedWp)

		// when
		err := svc.HandleUpdateWateringPlan(context.Background(), &event)

		// then
		assert.NoError(t, err)
		clusterRepo.AssertNotCalled(t, "Update")
	})
}

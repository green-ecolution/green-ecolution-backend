package treecluster

import (
	"context"
	"testing"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	storageMock "github.com/green-ecolution/green-ecolution-backend/internal/storage/_mock"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/green-ecolution/green-ecolution-backend/internal/worker"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestTreeClusterService_HandleUpdateTree(t *testing.T) {
	t.Run("should update tree cluster lat long and region and send treecluster update event", func(t *testing.T) {
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		eventManager := worker.NewEventManager(entities.EventTypeUpdateTreeCluster)
		svc := NewTreeClusterService(clusterRepo, treeRepo, regionRepo, eventManager)

		// event
		_, ch, _ := eventManager.Subscribe(entities.EventTypeUpdateTreeCluster)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go eventManager.Run(ctx)

		prevTc := entities.TreeCluster{
			ID: 1,
			Region: &entities.Region{
				ID:   1,
				Name: "Sandberg",
			},
			Latitude:  utils.P(54.776366336440255),
			Longitude: utils.P(9.451084144617182),
		}
		prevTree := entities.Tree{
			ID:          1,
			TreeCluster: &prevTc,
			Number:      "T001",
			Latitude:    54.776366336440255,
			Longitude:   9.451084144617182,
		}

		updatedTree := entities.Tree{
			ID:          1,
			TreeCluster: &prevTc,
			Number:      "T001",
			Latitude:    54.811733806341856,
			Longitude:   9.482958846410169,
		}

		updatedTc := entities.TreeCluster{
			ID: 1,
			Region: &entities.Region{
				ID:   2,
				Name: "Mürwik",
			},
			Latitude:  utils.P(54.811733806341856),
			Longitude: utils.P(9.482958846410169),
		}

		event := entities.NewEventUpdateTree(&prevTree, &updatedTree)

		clusterRepo.EXPECT().Update(mock.Anything, int32(1), mock.Anything).Return(nil)
		clusterRepo.EXPECT().GetByID(mock.Anything, int32(1)).Return(&updatedTc, nil)

		// when
		err := svc.HandleUpdateTree(context.Background(), &event)

		// then
		assert.NoError(t, err)
		select {
		case recievedEvent, ok := <-ch:
			assert.True(t, ok)
			e := recievedEvent.(entities.EventUpdateTreeCluster)
			assert.Equal(t, e.Prev, &prevTc)
			assert.Equal(t, e.New, &updatedTc)
		case <-time.After(1 * time.Second):
			t.Fatal("event was not received")
		}
	})

	t.Run("should not update tree cluster if treeclusters in event are nil", func(t *testing.T) {
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		eventManager := worker.NewEventManager(entities.EventTypeUpdateTreeCluster)
		svc := NewTreeClusterService(clusterRepo, treeRepo, regionRepo, eventManager)

		// event
		_, ch, _ := eventManager.Subscribe(entities.EventTypeUpdateTreeCluster)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go eventManager.Run(ctx)

		prevTree := entities.Tree{
			ID:          1,
			TreeCluster: nil,
			Number:      "T001",
			Latitude:    54.776366336440255,
			Longitude:   9.451084144617182,
		}

		updatedTree := entities.Tree{
			ID:          1,
			TreeCluster: nil,
			Number:      "T002",
			Latitude:    54.776366336440255,
			Longitude:   9.451084144617182,
		}

		event := entities.NewEventUpdateTree(&prevTree, &updatedTree)

		// when
		err := svc.HandleUpdateTree(context.Background(), &event)

		// then
		assert.NoError(t, err)
		clusterRepo.AssertNotCalled(t, "Update")
		clusterRepo.AssertNotCalled(t, "GetCenterPoint")
		regionRepo.AssertNotCalled(t, "GetByPoint")

		select {
		case <-ch:
			t.Fatalf("event was triggered but shouldn't have been")
		case <-time.After(100 * time.Millisecond):
			assert.True(t, true)
		}
	})

	t.Run("should not update tree cluster if tree has not changed location", func(t *testing.T) {
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		eventManager := worker.NewEventManager(entities.EventTypeUpdateTreeCluster)
		svc := NewTreeClusterService(clusterRepo, treeRepo, regionRepo, eventManager)

		// event
		_, ch, _ := eventManager.Subscribe(entities.EventTypeUpdateTreeCluster)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go eventManager.Run(ctx)

		tc := entities.TreeCluster{
			ID: 1,
			Region: &entities.Region{
				ID:   1,
				Name: "Sandberg",
			},
			Latitude:  utils.P(54.776366336440255),
			Longitude: utils.P(9.451084144617182),
		}
		prevTree := entities.Tree{
			ID:          1,
			TreeCluster: &tc,
			Number:      "T001",
			Latitude:    54.776366336440255,
			Longitude:   9.451084144617182,
		}

		updatedTree := entities.Tree{
			ID:          1,
			TreeCluster: &tc,
			Number:      "T002",
			Latitude:    54.776366336440255,
			Longitude:   9.451084144617182,
		}

		event := entities.NewEventUpdateTree(&prevTree, &updatedTree)

		// when
		err := svc.HandleUpdateTree(context.Background(), &event)

		// then
		assert.NoError(t, err)
		clusterRepo.AssertNotCalled(t, "Update")
		clusterRepo.AssertNotCalled(t, "GetCenterPoint")
		regionRepo.AssertNotCalled(t, "GetByPoint")

		select {
		case <-ch:
			t.Fatalf("event was triggered but shouldn't have been")
		case <-time.After(100 * time.Millisecond):
			assert.True(t, true)
		}
	})

	t.Run("should update if tree location is equal but tree has changed treecluster", func(t *testing.T) {
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		eventManager := worker.NewEventManager(entities.EventTypeUpdateTreeCluster)
		svc := NewTreeClusterService(clusterRepo, treeRepo, regionRepo, eventManager)

		// event
		_, ch, _ := eventManager.Subscribe(entities.EventTypeUpdateTreeCluster)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go eventManager.Run(ctx)

		prevTc := entities.TreeCluster{
			ID: 1,
			Region: &entities.Region{
				ID:   1,
				Name: "Sandberg",
			},
			Latitude:  utils.P(54.776366336440255),
			Longitude: utils.P(9.451084144617182),
		}
		prevTree := entities.Tree{
			ID:          1,
			TreeCluster: &prevTc,
			Number:      "T001",
			Latitude:    54.776366336440255,
			Longitude:   9.451084144617182,
		}

		newTc := entities.TreeCluster{
			ID: 2,
			Region: &entities.Region{
				ID:   1,
				Name: "Sandberg",
			},
			Latitude:  utils.P(54.776366336440255),
			Longitude: utils.P(9.451084144617182),
		}
		updatedTree := entities.Tree{
			ID:          1,
			TreeCluster: &newTc,
			Number:      "T002",
			Latitude:    54.776366336440255,
			Longitude:   9.451084144617182,
		}

		event := entities.NewEventUpdateTree(&prevTree, &updatedTree)

		clusterRepo.EXPECT().Update(mock.Anything, int32(1), mock.Anything).Return(nil)
		clusterRepo.EXPECT().Update(mock.Anything, int32(2), mock.Anything).Return(nil)
		clusterRepo.EXPECT().GetByID(mock.Anything, int32(1)).Return(&prevTc, nil)
		clusterRepo.EXPECT().GetByID(mock.Anything, int32(2)).Return(&newTc, nil)

		// when
		err := svc.HandleUpdateTree(context.Background(), &event)

		// then
		assert.NoError(t, err)
		clusterRepo.AssertNotCalled(t, "Update")
		clusterRepo.AssertNotCalled(t, "GetCenterPoint")
		regionRepo.AssertNotCalled(t, "GetByPoint")

		select {
		case _, ok := <-ch:
			assert.True(t, ok)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("event was not received")
		}
	})
}

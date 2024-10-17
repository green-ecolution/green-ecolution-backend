package treecluster

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	storageMock "github.com/green-ecolution/green-ecolution-backend/internal/storage/_mock"
	"github.com/stretchr/testify/assert"
)

func TestTreeClusterService_GetAll(t *testing.T) {
	ctx := context.Background()

	t.Run("should return all tree clusters when successful", func(t *testing.T) {
		// given
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := NewLocationUpdate(clusterRepo, treeRepo, regionRepo)
		svc := NewTreeClusterService(clusterRepo, treeRepo, regionRepo, locator)

		now := time.Now()
		expectedTrees := []*entities.Tree{
			{
				ID:           1,
				CreatedAt:    now,
				UpdatedAt:    now,
				Species:      "Oak",
				Number:       "T001",
				Latitude:     9.446741,
				Longitude:    54.801539,
				Description:  "A mature oak tree",
				PlantingYear: 2023,
				Readonly:     true,
			},
			{
				ID:           2,
				CreatedAt:    now,
				UpdatedAt:    now,
				Species:      "Pine",
				Number:       "T002",
				Latitude:     9.446700,
				Longitude:    54.801510,
				Description:  "A young pine tree",
				PlantingYear: 2023,
				Readonly:     true,
			},
		}

		expectedClusters := []*entities.TreeCluster{
			{
				ID:             1,
				CreatedAt:      now,
				UpdatedAt:      now,
				Name:           "Cluster 1",
				Address:        "123 Main St",
				Description:    "Test description",
				SoilCondition:  entities.TreeSoilConditionLehmig,
				Archived:       false,
				Latitude:       nil,
				Longitude:      nil,
				Trees:         expectedTrees,
			},
			{
				ID:             2,
				CreatedAt:      now,
				UpdatedAt:      now,
				Name:           "Cluster 2",
				Address:        "456 Second St",
				Description:    "Test description",
				SoilCondition:  entities.TreeSoilConditionSandig,
				Archived:       false,
				Latitude:       nil,
				Longitude:      nil,
				Trees:         []*entities.Tree{},
			},
		}

		clusterRepo.EXPECT().GetAll(ctx).Return(expectedClusters, nil)

		// when
		clusters, err := svc.GetAll(ctx)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedClusters, clusters)
	})

	t.Run("should return empty slice when no clusters are found", func(t *testing.T) {
		// given
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := NewLocationUpdate(clusterRepo, treeRepo, regionRepo)
		svc := NewTreeClusterService(clusterRepo, treeRepo, regionRepo, locator)

		clusterRepo.EXPECT().GetAll(ctx).Return([]*entities.TreeCluster{}, nil)

		// when
		clusters, err := svc.GetAll(ctx)

		// then
		assert.NoError(t, err)
		assert.Empty(t, clusters)
	})

	t.Run("should return error when GetAll fails", func(t *testing.T) {
		// given
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := NewLocationUpdate(clusterRepo, treeRepo, regionRepo)
		svc := NewTreeClusterService(clusterRepo, treeRepo, regionRepo, locator)

		expectedError := errors.New("GetAll failed")

		clusterRepo.EXPECT().GetAll(ctx).Return(nil, expectedError)

		// when
		clusters, err := svc.GetAll(ctx)

		// then
		assert.Error(t, err)
		assert.Nil(t, clusters)
		assert.EqualError(t, err, "500: " + expectedError.Error())
	})
}

func TestReady(t *testing.T) {
	t.Run("should return true if the service is ready", func(t *testing.T) {
		// given
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := NewLocationUpdate(clusterRepo, treeRepo, regionRepo)
		svc := NewTreeClusterService(clusterRepo, treeRepo, regionRepo, locator)

		// when
		ready := svc.Ready()

		// then
		assert.True(t, ready)
	})

	t.Run("should return false if the service is not ready", func(t *testing.T) {
		// given
		svc := NewTreeClusterService(nil, nil, nil, nil)

		// when
		ready := svc.Ready()

		// then
		assert.False(t, ready)
	})
}

package treecluster

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	service "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	storageMock "github.com/green-ecolution/green-ecolution-backend/internal/storage/_mock"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestTreeClusterService_GetAll(t *testing.T) {
	ctx := context.Background()

	t.Run("should return all tree clusters when successful", func(t *testing.T) {
		// given
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := service.NewMockGeoClusterLocator(t)
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
				Longitude:    9.446700,
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
				Latitude:       float64Ptr(9.446741),
				Longitude:      float64Ptr(54.801539),
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
		locator := service.NewMockGeoClusterLocator(t)
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
		locator := service.NewMockGeoClusterLocator(t)
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

func TestTreeClusterService_GetByID(t *testing.T) {
	ctx := context.Background()

	clusterRepo := storageMock.NewMockTreeClusterRepository(t)
	treeRepo := storageMock.NewMockTreeRepository(t)
	regionRepo := storageMock.NewMockRegionRepository(t)
	locator := service.NewMockGeoClusterLocator(t)
	svc := NewTreeClusterService(clusterRepo, treeRepo, regionRepo, locator)

	t.Run("should return tree cluster when found", func(t *testing.T) {
		id := int32(1)
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

		expectedCluster := &entities.TreeCluster{
			ID:             id,
			CreatedAt:      now,
			UpdatedAt:      now,
			Name:           "Cluster 1",
			Address:        "123 Main St",
			Description:    "Test description",
			SoilCondition:  entities.TreeSoilConditionLehmig,
			Archived:       false,
			Latitude:       float64Ptr(9.446741),
			Longitude:      float64Ptr(54.801539),
			Trees:         expectedTrees,
		}

		clusterRepo.EXPECT().GetByID(ctx, id).Return(expectedCluster, nil)

		// when
		cluster, err := svc.GetByID(ctx, id)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedCluster, cluster)
	})

	t.Run("should return error if tree cluster not found", func(t *testing.T) {
		id := int32(2)

		expectedError := storage.ErrEntityNotFound
		clusterRepo.EXPECT().GetByID(ctx, id).Return(nil, expectedError)

		// when
		cluster, err := svc.GetByID(ctx, id)

		// then
		assert.Error(t, err)
		assert.Nil(t, cluster)
		assert.EqualError(t, err, "404: " + expectedError.Error())
	})

	t.Run("should return error for unexpected repository error", func(t *testing.T) {
		id := int32(3)
		expectedError := errors.New("unexpected error")

		// Set expectation for GetByID
		clusterRepo.EXPECT().GetByID(ctx, id).Return(nil, expectedError)

		// when
		cluster, err := svc.GetByID(ctx, id)

		// then
		assert.Error(t, err)
		assert.Nil(t, cluster)
		assert.EqualError(t, err, "500: " + expectedError.Error())
	})
}

func TestTreeClusterService_Create(t *testing.T) {
	t.Run("should successfully create a new tree cluster", func(t *testing.T) {
		// Given
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := service.NewMockGeoClusterLocator(t)
		svc := NewTreeClusterService(clusterRepo, treeRepo, regionRepo, locator)

		now := time.Now()
		tc := &entities.TreeClusterCreate{
			Name:        "New Cluster",
			Address:     "123 Forest Lane",
			Description: "A test cluster",
			SoilCondition:  entities.TreeSoilConditionLehmig,
			TreeIDs:     []*int32{ptrToInt32(1), ptrToInt32(2)},
		}

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

		expectedCluster := &entities.TreeCluster{
			ID:             10,
			CreatedAt:      now,
			UpdatedAt:      now,
			Name:           "Cluster 1",
			Address:        "123 Main St",
			Description:    "Test description",
			SoilCondition:  entities.TreeSoilConditionLehmig,
			Archived:       false,
			Latitude:       float64Ptr(9.446741),
			Longitude:      float64Ptr(54.801539),
			Trees:         expectedTrees,
		}

		treeRepo.EXPECT().GetTreesByIDs(context.Background(), []int32{1, 2}).Return(expectedTrees, nil)
		clusterRepo.EXPECT().Create(context.Background(), mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(expectedCluster, nil)
		locator.EXPECT().UpdateCluster(context.Background(), &expectedCluster.ID).Return(nil)

		// When
		result, err := svc.Create(context.Background(), tc)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, expectedCluster, result)
	})
}

func TestTreeClusterService_Delete(t *testing.T) {
	ctx := context.Background()

	clusterRepo := storageMock.NewMockTreeClusterRepository(t)
	treeRepo := storageMock.NewMockTreeRepository(t)
	regionRepo := storageMock.NewMockRegionRepository(t)
	locator := service.NewMockGeoClusterLocator(t)
	svc := NewTreeClusterService(clusterRepo, treeRepo, regionRepo, locator)

    t.Run("should successfully delete a tree cluster", func(t *testing.T) {
        id := int32(1)

        clusterRepo.EXPECT().GetByID(ctx, id).Return(&entities.TreeCluster{}, nil)
        treeRepo.EXPECT().UnlinkTreeClusterID(ctx, id).Return(nil)
        clusterRepo.EXPECT().Delete(ctx, id).Return(nil)

		// when
        err := svc.Delete(ctx, id)

		// then
        assert.NoError(t, err)
    })

	t.Run("should return error if tree cluster not found", func(t *testing.T) {
		id := int32(2)

		expectedError := storage.ErrEntityNotFound
		clusterRepo.EXPECT().GetByID(ctx, id).Return(nil, expectedError)

		// when
        err := svc.Delete(ctx, id)
		
		// then
		assert.Error(t, err)
		assert.EqualError(t, err, "404: " + expectedError.Error())
	})

	t.Run("should return error if unlinking tree cluster ID fails", func(t *testing.T) {
        id := int32(3)
        expectedError := errors.New("failed to unlink")

        clusterRepo.EXPECT().GetByID(ctx, id).Return(&entities.TreeCluster{}, nil)
        treeRepo.EXPECT().UnlinkTreeClusterID(ctx, id).Return(expectedError)

        // when
        err := svc.Delete(ctx, id)

		// then
        assert.Error(t, err)
		assert.EqualError(t, err, "500: " + expectedError.Error())
    })

	t.Run("should return error if deleting tree cluster fails", func(t *testing.T) {
        id := int32(4)
        expectedError := errors.New("failed to delete")

        // Set expectations
        clusterRepo.EXPECT().GetByID(ctx, id).Return(&entities.TreeCluster{}, nil)
        treeRepo.EXPECT().UnlinkTreeClusterID(ctx, id).Return(nil)
		clusterRepo.EXPECT().Delete(ctx, id).Return(expectedError)

        // when
        err := svc.Delete(ctx, id)

        // then
		assert.Error(t, err)
		assert.EqualError(t, err, "500: " + expectedError.Error())
    })
}

func TestReady(t *testing.T) {
	t.Run("should return true if the service is ready", func(t *testing.T) {
		// given
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := service.NewMockGeoClusterLocator(t)
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

func float64Ptr(f float64) *float64 {
	return &f
}

func ptrToInt32(i int32) *int32 {
	return &i
}

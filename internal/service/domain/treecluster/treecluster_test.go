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
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := service.NewMockGeoClusterLocator(t)
		svc := NewTreeClusterService(clusterRepo, treeRepo, regionRepo, locator)

		expectedClusters := getTestTreeClusters()
		clusterRepo.EXPECT().GetAll(ctx).Return(expectedClusters, nil)

		// when
		clusters, err := svc.GetAll(ctx)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedClusters, clusters)
	})

	t.Run("should return empty slice when no clusters are found", func(t *testing.T) {
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
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := service.NewMockGeoClusterLocator(t)
		svc := NewTreeClusterService(clusterRepo, treeRepo, regionRepo, locator)

		expectedErr := errors.New("GetAll failed")

		clusterRepo.EXPECT().GetAll(ctx).Return(nil, expectedErr)

		// when
		clusters, err := svc.GetAll(ctx)

		// then
		assert.Error(t, err)
		assert.Nil(t, clusters)
		assert.EqualError(t, err, handleError(expectedErr).Error())
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
		expectedCluster := getTestTreeClusters()[0]
		clusterRepo.EXPECT().GetByID(ctx, id).Return(expectedCluster, nil)

		// when
		cluster, err := svc.GetByID(ctx, id)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedCluster, cluster)
	})

	t.Run("should return error if tree cluster not found", func(t *testing.T) {
		id := int32(2)
		expectedErr := storage.ErrEntityNotFound
		clusterRepo.EXPECT().GetByID(ctx, id).Return(nil, expectedErr)

		// when
		cluster, err := svc.GetByID(ctx, id)

		// then
		assert.Error(t, err)
		assert.Nil(t, cluster)
		assert.EqualError(t, err, handleError(expectedErr).Error())
	})

	t.Run("should return error for unexpected repository error", func(t *testing.T) {
		id := int32(3)
		expectedErr := errors.New("unexpected error")

		clusterRepo.EXPECT().GetByID(ctx, id).Return(nil, expectedErr)

		// when
		cluster, err := svc.GetByID(ctx, id)

		// then
		assert.Error(t, err)
		assert.Nil(t, cluster)
		assert.EqualError(t, err, handleError(expectedErr).Error())
	})
}

func TestTreeClusterService_Create(t *testing.T) {
	ctx := context.Background()
	newCluster := &entities.TreeClusterCreate{
		Name:          "Cluster 1",
		Address:       "123 Main St",
		Description:   "Test description",
		SoilCondition:  entities.TreeSoilConditionLehmig,
		TreeIDs:     []*int32{ptrToInt32(1), ptrToInt32(2)},
	}

	t.Run("should successfully create a new tree cluster", func(t *testing.T) {
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := service.NewMockGeoClusterLocator(t)
		svc := NewTreeClusterService(clusterRepo, treeRepo, regionRepo, locator)

		expectedCluster := getTestTreeClusters()[0]
		expectedTrees := getTestTrees()

		treeRepo.EXPECT().GetTreesByIDs(
			ctx,
			[]int32{1, 2},
		).Return(expectedTrees, nil)

		clusterRepo.EXPECT().Create(
			ctx,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(expectedCluster, nil)

		locator.EXPECT().UpdateCluster(
			ctx,
			&expectedCluster.ID,
		).Return(nil)

		// when
		result, err := svc.Create(ctx, newCluster)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedCluster, result)
	})

	t.Run("should successfully create a new tree cluster with empty trees", func(t *testing.T) {
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := service.NewMockGeoClusterLocator(t)
		svc := NewTreeClusterService(clusterRepo, treeRepo, regionRepo, locator)
	
		newCluster := &entities.TreeClusterCreate{
			Name:          "Cluster 1",
			Address:       "123 Main St",
			Description:   "Test description",
			SoilCondition: entities.TreeSoilConditionLehmig,
			TreeIDs:       []*int32{},
		}

		expectedCluster := getTestTreeClusters()[1]

		treeRepo.EXPECT().GetTreesByIDs(
			ctx,
			[]int32{},
		).Return(nil, nil)

		clusterRepo.EXPECT().Create(
			ctx,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(expectedCluster, nil)

		locator.EXPECT().UpdateCluster(
			ctx,
			&expectedCluster.ID,
		).Return(nil)
	
		// when
		result, err := svc.Create(ctx, newCluster)
	
		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedCluster, result)
	})

	t.Run("should return an error when getting trees fails", func(t *testing.T) {
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := service.NewMockGeoClusterLocator(t)
		svc := NewTreeClusterService(clusterRepo, treeRepo, regionRepo, locator)

		expectedErr := storage.ErrTreeNotFound

		treeRepo.EXPECT().GetTreesByIDs(
			ctx, 
			[]int32{1, 2},
		).Return(nil, expectedErr)
	
		// when
		result, err := svc.Create(ctx, newCluster)
	
		// then
		assert.Nil(t, result)
		assert.EqualError(t, err, handleError(expectedErr).Error())
	})

	t.Run("should return an error when creating cluster fails", func(t *testing.T) {
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := service.NewMockGeoClusterLocator(t)
		svc := NewTreeClusterService(clusterRepo, treeRepo, regionRepo, locator)
	
		expectedErr := errors.New("Failed to create cluster")
		expectedTrees := getTestTrees()
	
		treeRepo.EXPECT().GetTreesByIDs(
			ctx, 
			[]int32{1, 2},
		).Return(expectedTrees, nil)

		clusterRepo.EXPECT().Create(
			ctx, 
			mock.Anything, 
			mock.Anything, 
			mock.Anything, 
			mock.Anything,
		).Return(nil, expectedErr)
	
		// when
		result, err := svc.Create(ctx, newCluster)
	
		// then
		assert.Nil(t, result)
		assert.EqualError(t, err, handleError(expectedErr).Error())
	})
}

func TestTreeClusterService_Update(t *testing.T) {
	ctx := context.Background()
	clusterID := int32(1)
	updatedCluster := &entities.TreeClusterUpdate{
		Name:          "Cluster 1",
		Address:       "123 Main St",
		Description:   "Test description",
		SoilCondition: entities.TreeSoilConditionLehmig,
		TreeIDs:     []*int32{ptrToInt32(1), ptrToInt32(2)},
	}

	t.Run("should successfully update a tree cluster", func(t *testing.T) {
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := service.NewMockGeoClusterLocator(t)
		svc := NewTreeClusterService(clusterRepo, treeRepo, regionRepo, locator)

		expectedCluster := getTestTreeClusters()[0]
		expectedTrees := getTestTrees()

		treeRepo.EXPECT().GetTreesByIDs(
			ctx,
			[]int32{1, 2},
		).Return(expectedTrees, nil)

		clusterRepo.EXPECT().Update(
			ctx,
			clusterID,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(expectedCluster, nil)

		locator.EXPECT().UpdateCluster(
			ctx,
			&clusterID,
		).Return(nil)

		// when
		result, err := svc.Update(ctx, clusterID, updatedCluster)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedCluster, result)
	})


	t.Run("should successfully update a tree cluster with empty tree IDs", func(t *testing.T) {
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := service.NewMockGeoClusterLocator(t)
		svc := NewTreeClusterService(clusterRepo, treeRepo, regionRepo, locator)

		updatedClusterEmptyTrees := &entities.TreeClusterUpdate{
			Name:          "Cluster 1",
			Address:       "123 Main St",
			Description:   "Test description",
			SoilCondition: entities.TreeSoilConditionLehmig,
			TreeIDs:       []*int32{},
		}

		expectedCluster := getTestTreeClusters()[1]

		treeRepo.EXPECT().GetTreesByIDs(
			ctx,
			[]int32{},
		).Return(nil, nil)

		clusterRepo.EXPECT().Update(
			ctx,
			clusterID,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(expectedCluster, nil)

		locator.EXPECT().UpdateCluster(
			ctx,
			&clusterID,
		).Return(nil)

		// when
		result, err := svc.Update(ctx, clusterID, updatedClusterEmptyTrees)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedCluster, result)
	})

	t.Run("should return an error when no trees are found", func(t *testing.T) {
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := service.NewMockGeoClusterLocator(t)
		svc := NewTreeClusterService(clusterRepo, treeRepo, regionRepo, locator)

		treeRepo.EXPECT().GetTreesByIDs(
			ctx, 
			[]int32{1, 2},
		).Return(nil, storage.ErrTreeNotFound)

		// when
		result, err := svc.Update(context.Background(), clusterID, updatedCluster)

		// then
		assert.Nil(t, result)
		assert.EqualError(t, err, handleError(storage.ErrTreeNotFound).Error())
	})

	t.Run("should return an error when the update fails", func(t *testing.T) {
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := service.NewMockGeoClusterLocator(t)
		svc := NewTreeClusterService(clusterRepo, treeRepo, regionRepo, locator)

		expectedErr := errors.New("failed to update cluster")
		expectedTrees := getTestTrees()

		treeRepo.EXPECT().GetTreesByIDs(
			ctx,
			[]int32{1,2},
		).Return(expectedTrees, nil)

		clusterRepo.EXPECT().Update(
			ctx,
			clusterID,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedErr)

		// when
		result, err := svc.Update(context.Background(), clusterID, updatedCluster)

		// then
		assert.Nil(t, result)
		assert.EqualError(t, err, handleError(expectedErr).Error())
	})

	t.Run("should return an error when cluster ID does not exist", func(t *testing.T) {
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := service.NewMockGeoClusterLocator(t)
		svc := NewTreeClusterService(clusterRepo, treeRepo, regionRepo, locator)

		expectedTrees := getTestTrees()

		treeRepo.EXPECT().GetTreesByIDs(
			ctx,
			[]int32{1,2},
		).Return(expectedTrees, nil)

		clusterRepo.EXPECT().Update(
			ctx,
			clusterID,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, storage.ErrEntityNotFound)

		// when
		result, err := svc.Update(ctx, clusterID, updatedCluster)

		// then
		assert.Nil(t, result)
		assert.EqualError(t, err, handleError(storage.ErrEntityNotFound).Error())
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

        clusterRepo.EXPECT().GetByID(ctx, id).Return(getTestTreeClusters()[0], nil)
        treeRepo.EXPECT().UnlinkTreeClusterID(ctx, id).Return(nil)
        clusterRepo.EXPECT().Delete(ctx, id).Return(nil)

		// when
        err := svc.Delete(ctx, id)

		// then
        assert.NoError(t, err)
    })

	t.Run("should return error if tree cluster not found", func(t *testing.T) {
		id := int32(2)

		expectedErr := storage.ErrEntityNotFound
		clusterRepo.EXPECT().GetByID(ctx, id).Return(nil, expectedErr)

		// when
        err := svc.Delete(ctx, id)
		
		// then
		assert.Error(t, err)
		assert.EqualError(t, err, handleError(expectedErr).Error())
	})

	t.Run("should return error if unlinking tree cluster ID fails", func(t *testing.T) {
        id := int32(3)
        expectedErr := errors.New("failed to unlink")

        clusterRepo.EXPECT().GetByID(ctx, id).Return(getTestTreeClusters()[0], nil)
        treeRepo.EXPECT().UnlinkTreeClusterID(ctx, id).Return(expectedErr)

        // when
        err := svc.Delete(ctx, id)

		// then
        assert.Error(t, err)
		assert.EqualError(t, err, handleError(expectedErr).Error())
    })

	t.Run("should return error if deleting tree cluster fails", func(t *testing.T) {
        id := int32(4)
        expectedErr := errors.New("failed to delete")

        clusterRepo.EXPECT().GetByID(ctx, id).Return(getTestTreeClusters()[0], nil)
        treeRepo.EXPECT().UnlinkTreeClusterID(ctx, id).Return(nil)
		clusterRepo.EXPECT().Delete(ctx, id).Return(expectedErr)

        // when
        err := svc.Delete(ctx, id)

        // then
		assert.Error(t, err)
		assert.EqualError(t, err, handleError(expectedErr).Error())
    })
}

func TestReady(t *testing.T) {
	t.Run("should return true if the service is ready", func(t *testing.T) {
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
		svc := NewTreeClusterService(nil, nil, nil, nil)

		// when
		ready := svc.Ready()

		// then
		assert.False(t, ready)
	})
}

func getTestTreeClusters() []*entities.TreeCluster {
	now := time.Now()

	return []*entities.TreeCluster{
		{
			ID:            1,
			CreatedAt:     now,
			UpdatedAt:     now,
			Name:          "Cluster 1",
			Address:       "123 Main St",
			Description:   "Test description",
			SoilCondition: entities.TreeSoilConditionLehmig,
			Archived:      false,
			Latitude:      float64Ptr(9.446741),
			Longitude:     float64Ptr(54.801539),
			Trees:         getTestTrees(),
		},
		{
			ID:            2,
			CreatedAt:     now,
			UpdatedAt:     now,
			Name:          "Cluster 2",
			Address:       "456 Second St",
			Description:   "Test description",
			SoilCondition: entities.TreeSoilConditionSandig,
			Archived:      false,
			Latitude:      nil,
			Longitude:     nil,
			Trees:        []*entities.Tree{},
		},
	}
}

func getTestTrees() []*entities.Tree {
	now := time.Now()

	return []*entities.Tree{
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
}

func float64Ptr(f float64) *float64 {
	return &f
}

func ptrToInt32(i int32) *int32 {
	return &i
}

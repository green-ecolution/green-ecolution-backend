package tree

import (
	"context"
	"errors"
	service "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
	"github.com/green-ecolution/green-ecolution-backend/internal/service/domain/treecluster"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	storageMock "github.com/green-ecolution/green-ecolution-backend/internal/storage/_mock"
	"github.com/stretchr/testify/assert"
)

func TestTreeService_GetAll(t *testing.T) {
	ctx := context.Background()

	t.Run("should return all trees when successful", func(t *testing.T) {
		// given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		imageRepo := storageMock.NewMockImageRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)

		locator := treecluster.NewLocationUpdate(clusterRepo, treeRepo, regionRepo)

		svc := NewTreeService(treeRepo, sensorRepo, imageRepo, clusterRepo, locator)

		expectedTrees := getTestTrees()
		treeRepo.EXPECT().GetAll(ctx).Return(expectedTrees, nil)

		// when
		trees, err := svc.GetAll(ctx)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedTrees, trees)
	})

	t.Run("should return empty slice when no trees are found", func(t *testing.T) {
		// given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		imageRepo := storageMock.NewMockImageRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)

		locator := treecluster.NewLocationUpdate(clusterRepo, treeRepo, regionRepo)

		svc := NewTreeService(treeRepo, sensorRepo, imageRepo, clusterRepo, locator)

		treeRepo.EXPECT().GetAll(ctx).Return([]*entities.Tree{}, nil)

		// when
		trees, err := svc.GetAll(ctx)

		// then
		assert.NoError(t, err)
		assert.Empty(t, trees)
	})

	t.Run("should return error when GetAll fails", func(t *testing.T) {
		// given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		imageRepo := storageMock.NewMockImageRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)

		locator := treecluster.NewLocationUpdate(clusterRepo, treeRepo, regionRepo)

		svc := NewTreeService(treeRepo, sensorRepo, imageRepo, clusterRepo, locator)

		expectedError := errors.New("GetAll failed")

		treeRepo.EXPECT().GetAll(ctx).Return(nil, expectedError)

		// when
		trees, err := svc.GetAll(ctx)

		// then
		assert.Error(t, err)
		assert.Nil(t, trees)
		assert.EqualError(t, err, "500: "+expectedError.Error())
	})
}

func TestTreeService_GetByID(t *testing.T) {
	ctx := context.Background()

	// Mocks
	treeRepo := storageMock.NewMockTreeRepository(t)
	sensorRepo := storageMock.NewMockSensorRepository(t)
	imageRepo := storageMock.NewMockImageRepository(t)
	clusterRepo := storageMock.NewMockTreeClusterRepository(t)
	regionRepo := storageMock.NewMockRegionRepository(t)

	locator := treecluster.NewLocationUpdate(clusterRepo, treeRepo, regionRepo)
	svc := NewTreeService(treeRepo, sensorRepo, imageRepo, clusterRepo, locator)

	t.Run("should return tree when found", func(t *testing.T) {
		id := int32(1)
		expectedTree := getTestTrees()[0]
		treeRepo.EXPECT().GetByID(ctx, id).Return(expectedTree, nil)

		// when
		tree, err := svc.GetByID(ctx, id)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedTree, tree)
	})

	t.Run("should return error if tree not found", func(t *testing.T) {
		id := int32(2)
		expectedError := storage.ErrEntityNotFound
		treeRepo.EXPECT().GetByID(ctx, id).Return(nil, expectedError)

		// when
		tree, err := svc.GetByID(ctx, id)

		// then
		assert.Error(t, err)
		assert.Nil(t, tree)
		assert.EqualError(t, err, "404: "+expectedError.Error())
	})

	t.Run("should return error for unexpected repository error", func(t *testing.T) {
		id := int32(3)
		expectedError := errors.New("unexpected error")

		// Set expectation for GetByID
		treeRepo.EXPECT().GetByID(ctx, id).Return(nil, expectedError)

		// when
		tree, err := svc.GetByID(ctx, id)

		// then
		assert.Error(t, err)
		assert.Nil(t, tree)
		assert.EqualError(t, err, "500: "+expectedError.Error())
	})
}

func TestTreeService_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully delete a tree", func(t *testing.T) {
		// Given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		imageRepo := storageMock.NewMockImageRepository(t)
		treeClusterRepo := storageMock.NewMockTreeClusterRepository(t)
		locator := service.NewMockGeoClusterLocator(t)

		svc := NewTreeService(treeRepo, sensorRepo, imageRepo, treeClusterRepo, locator)

		expectedTree := getTestTrees()[0]
		expectedTree.TreeCluster = getTestTreeClusters()[0]

		// Mock expectations
		treeRepo.EXPECT().GetByID(ctx, expectedTree.ID).Return(expectedTree, nil)
		treeRepo.EXPECT().Delete(ctx, expectedTree.ID).Return(nil)
		locator.EXPECT().UpdateCluster(ctx, &expectedTree.TreeCluster.ID).Return(nil)

		// When
		err := svc.Delete(ctx, expectedTree.ID)

		// Then
		assert.NoError(t, err)
	})

	t.Run("should return error if tree not found", func(t *testing.T) {
		// Given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		imageRepo := storageMock.NewMockImageRepository(t)
		treeClusterRepo := storageMock.NewMockTreeClusterRepository(t)
		locator := service.NewMockGeoClusterLocator(t)

		svc := NewTreeService(treeRepo, sensorRepo, imageRepo, treeClusterRepo, locator)

		id := int32(1)
		expectedError := storage.ErrTreeNotFound

		// Mock expectations
		treeRepo.EXPECT().GetByID(ctx, id).Return(nil, expectedError)

		// When
		err := svc.Delete(ctx, id)

		// Then
		assert.Error(t, err)
		assert.EqualError(t, err, handleError(expectedError).Error())
	})
	t.Run("should return error if tree deletion fails", func(t *testing.T) {
		// Given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		imageRepo := storageMock.NewMockImageRepository(t)
		treeClusterRepo := storageMock.NewMockTreeClusterRepository(t)
		locator := service.NewMockGeoClusterLocator(t)

		svc := NewTreeService(treeRepo, sensorRepo, imageRepo, treeClusterRepo, locator)

		expectedTree := getTestTrees()[0]
		expectedTree.TreeCluster = getTestTreeClusters()[0]
		expectedError := errors.New("deletion failed")

		// Mock expectations
		treeRepo.EXPECT().GetByID(ctx, expectedTree.ID).Return(expectedTree, nil)
		treeRepo.EXPECT().Delete(ctx, expectedTree.ID).Return(expectedError)

		// When
		err := svc.Delete(ctx, expectedTree.ID)

		// Then
		assert.Error(t, err)
		assert.EqualError(t, err, handleError(expectedError).Error())
	})
	t.Run("should return error if updating cluster fails after deleting tree", func(t *testing.T) {
		// Given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		imageRepo := storageMock.NewMockImageRepository(t)
		treeClusterRepo := storageMock.NewMockTreeClusterRepository(t)
		locator := service.NewMockGeoClusterLocator(t)

		svc := NewTreeService(treeRepo, sensorRepo, imageRepo, treeClusterRepo, locator)

		expectedTree := getTestTrees()[0]
		expectedTree.TreeCluster = getTestTreeClusters()[0]
		expectedError := errors.New("cluster update failed")

		// Mock expectations
		treeRepo.EXPECT().GetByID(ctx, expectedTree.ID).Return(expectedTree, nil)
		treeRepo.EXPECT().Delete(ctx, expectedTree.ID).Return(nil)
		locator.EXPECT().UpdateCluster(ctx, &expectedTree.TreeCluster.ID).Return(expectedError)

		// When
		err := svc.Delete(ctx, expectedTree.ID)

		// Then
		assert.Error(t, err)
		assert.EqualError(t, err, handleError(expectedError).Error())
	})
	t.Run("should delete a tree without triggering cluster update when tree has no cluster", func(t *testing.T) {
		// Given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		imageRepo := storageMock.NewMockImageRepository(t)
		treeClusterRepo := storageMock.NewMockTreeClusterRepository(t)
		locator := service.NewMockGeoClusterLocator(t)

		svc := NewTreeService(treeRepo, sensorRepo, imageRepo, treeClusterRepo, locator)

		expectedTree := getTestTrees()[0]
		expectedTree.TreeCluster = nil // Tree has no cluster

		// Mock expectations
		treeRepo.EXPECT().GetByID(ctx, expectedTree.ID).Return(expectedTree, nil)
		treeRepo.EXPECT().Delete(ctx, expectedTree.ID).Return(nil)

		// Ensure UpdateCluster is not called
		locator.AssertNotCalled(t, "UpdateCluster", ctx, mock.Anything)

		// When
		err := svc.Delete(ctx, expectedTree.ID)

		// Then
		assert.NoError(t, err)
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
			Trees:         []*entities.Tree{},
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

func getTestSensors() []*entities.Sensor {
	now := time.Now()

	return []*entities.Sensor{
		{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
			Status:    entities.SensorStatusUnknown,
			Data:      nil,
		},
		{
			ID:        2,
			CreatedAt: now,
			UpdatedAt: now,
			Status:    entities.SensorStatusUnknown,
			Data:      nil,
		},
	}

}

func float64Ptr(f float64) *float64 {
	return &f
}

func ptrToInt32(i int32) *int32 {
	return &i
}

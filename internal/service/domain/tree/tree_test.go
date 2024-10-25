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
func TestTreeService_Ready(t *testing.T) {
	t.Run("should return true when treeRepo and sensorRepo are initialized", func(t *testing.T) {
		// Given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)

		svc := NewTreeService(treeRepo, sensorRepo, nil, nil, nil)

		// When
		result := svc.Ready()

		// Then
		assert.True(t, result)
	})

	t.Run("should return false when treeRepo is nil", func(t *testing.T) {
		// Given
		sensorRepo := storageMock.NewMockSensorRepository(t)

		svc := NewTreeService(nil, sensorRepo, nil, nil, nil)

		// When
		result := svc.Ready()

		// Then
		assert.False(t, result)
	})

	t.Run("should return false when sensorRepo is nil", func(t *testing.T) {
		// Given
		treeRepo := storageMock.NewMockTreeRepository(t)

		svc := NewTreeService(treeRepo, nil, nil, nil, nil)

		// When
		result := svc.Ready()

		// Then
		assert.False(t, result)
	})

	t.Run("should return false when both treeRepo and sensorRepo are nil", func(t *testing.T) {
		// Given
		svc := NewTreeService(nil, nil, nil, nil, nil)

		// When
		result := svc.Ready()

		// Then
		assert.False(t, result)
	})
}
func TestTreeService_Create(t *testing.T) {
	ctx := context.Background()
	newTreeCreate := &entities.TreeCreate{
		Species:       "Oak",
		Latitude:      54.801539,
		Longitude:     9.446741,
		PlantingYear:  2023,
		Number:        "T001",
		TreeClusterID: ptrToInt32(1),
		SensorID:      ptrToInt32(1),
	}
	t.Run("should successfully create a new tree", func(t *testing.T) {
		// Given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		imageRepo := storageMock.NewMockImageRepository(t)
		treeClusterRepo := storageMock.NewMockTreeClusterRepository(t)
		locator := service.NewMockGeoClusterLocator(t)

		svc := NewTreeService(treeRepo, sensorRepo, imageRepo, treeClusterRepo, locator)

		expectedTree := getTestTrees()[0]
		expectedCluster := getTestTreeClusters()[0]
		expectedSensor := getTestSensors()[0]

		// Mock expectations
		treeClusterRepo.EXPECT().GetByID(ctx, int32(1)).Return(expectedCluster, nil)
		sensorRepo.EXPECT().GetByID(ctx, int32(1)).Return(expectedSensor, nil)
		treeRepo.EXPECT().Create(ctx,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything).Return(expectedTree, nil)
		locator.EXPECT().UpdateCluster(ctx, newTreeCreate.TreeClusterID).Return(nil)

		// When
		result, err := svc.Create(ctx, newTreeCreate)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, expectedTree, result)
	})

	t.Run("should return validation error", func(t *testing.T) {
		// Given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		imageRepo := storageMock.NewMockImageRepository(t)
		treeClusterRepo := storageMock.NewMockTreeClusterRepository(t)
		locator := service.NewMockGeoClusterLocator(t)

		svc := NewTreeService(treeRepo, sensorRepo, imageRepo, treeClusterRepo, locator)

		invalidTreeCreate := &entities.TreeCreate{
			Species:      "Oak",
			Latitude:     0,  // Invalid: must be between -90 and 90
			Longitude:    0,  // Invalid: must be between -180 and 180
			PlantingYear: 0,  // Invalid: PlantingYear is required
			Number:       "", // Invalid: Number is required
		}

		// When
		result, err := svc.Create(ctx, invalidTreeCreate)

		// Then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validation error")
	})
	t.Run("should return error when fetching TreeCluster fails", func(t *testing.T) {
		// Given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		imageRepo := storageMock.NewMockImageRepository(t)
		treeClusterRepo := storageMock.NewMockTreeClusterRepository(t)
		locator := service.NewMockGeoClusterLocator(t)

		svc := NewTreeService(treeRepo, sensorRepo, imageRepo, treeClusterRepo, locator)

		expectedError := storage.ErrTreeClusterNotFound

		// Mock expectations
		treeClusterRepo.EXPECT().GetByID(ctx, *newTreeCreate.TreeClusterID).Return(nil, expectedError)

		// When
		result, err := svc.Create(ctx, newTreeCreate)

		// Then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, handleError(expectedError).Error())
	})

	t.Run("should return error when fetching Sensor fails", func(t *testing.T) {
		// Given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		imageRepo := storageMock.NewMockImageRepository(t)
		treeClusterRepo := storageMock.NewMockTreeClusterRepository(t)
		locator := service.NewMockGeoClusterLocator(t)

		svc := NewTreeService(treeRepo, sensorRepo, imageRepo, treeClusterRepo, locator)

		expectedError := storage.ErrSensorNotFound
		expectedCluster := getTestTreeClusters()[0]

		// Mock expectations
		treeClusterRepo.EXPECT().GetByID(ctx, int32(1)).Return(expectedCluster, nil)
		sensorRepo.EXPECT().GetByID(ctx, *newTreeCreate.SensorID).Return(nil, expectedError)

		// When
		result, err := svc.Create(ctx, newTreeCreate)

		// Then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, handleError(expectedError).Error())
	})

	t.Run("should return error when creating tree fails", func(t *testing.T) {
		// Given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		imageRepo := storageMock.NewMockImageRepository(t)
		treeClusterRepo := storageMock.NewMockTreeClusterRepository(t)
		locator := service.NewMockGeoClusterLocator(t)

		svc := NewTreeService(treeRepo, sensorRepo, imageRepo, treeClusterRepo, locator)

		expectedCluster := getTestTreeClusters()[0]
		expectedSensor := getTestSensors()[0]
		expectedError := errors.New("tree creation failed")

		// Mock expectations
		treeClusterRepo.EXPECT().GetByID(ctx, *newTreeCreate.TreeClusterID).Return(expectedCluster, nil)
		sensorRepo.EXPECT().GetByID(ctx, *newTreeCreate.SensorID).Return(expectedSensor, nil)
		treeRepo.EXPECT().Create(ctx,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything).Return(nil, expectedError)

		// When
		result, err := svc.Create(ctx, newTreeCreate)

		// Then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, handleError(expectedError).Error())
	})

	t.Run("should return error when updating cluster fails", func(t *testing.T) {
		// Given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		imageRepo := storageMock.NewMockImageRepository(t)
		treeClusterRepo := storageMock.NewMockTreeClusterRepository(t)
		locator := service.NewMockGeoClusterLocator(t)

		svc := NewTreeService(treeRepo, sensorRepo, imageRepo, treeClusterRepo, locator)

		expectedTree := getTestTrees()[0]
		expectedCluster := getTestTreeClusters()[0]
		expectedSensor := getTestSensors()[0]
		expectedError := errors.New("cluster update failed")

		// Mock expectations
		treeClusterRepo.EXPECT().GetByID(ctx, *newTreeCreate.TreeClusterID).Return(expectedCluster, nil)
		sensorRepo.EXPECT().GetByID(ctx, *newTreeCreate.SensorID).Return(expectedSensor, nil)
		treeRepo.EXPECT().Create(ctx,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything).Return(expectedTree, nil)
		locator.EXPECT().UpdateCluster(ctx, newTreeCreate.TreeClusterID).Return(expectedError)

		// When
		result, err := svc.Create(ctx, newTreeCreate)

		// Then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, handleError(expectedError).Error())
	})
}

func TestTreeService_Update(t *testing.T) {
	ctx := context.Background()

	id := int32(1)
	treeUpdate := &entities.TreeUpdate{
		TreeClusterID: ptrToInt32(1),
		SensorID:      ptrToInt32(1),
		PlantingYear:  2023,
		Species:       "Oak",
		Number:        "T001",
		Latitude:      54.801539,
		Longitude:     9.446741,
		Description:   "Updated description",
	}

	updatedTree := getTestTrees()[0]
	updatedTree.Description = treeUpdate.Description
	updatedTree.TreeCluster = getTestTreeClusters()[1]

	t.Run("should successfully update a tree", func(t *testing.T) {
		// Given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		imageRepo := storageMock.NewMockImageRepository(t)
		treeClusterRepo := storageMock.NewMockTreeClusterRepository(t)
		locator := service.NewMockGeoClusterLocator(t)

		svc := NewTreeService(treeRepo, sensorRepo, imageRepo, treeClusterRepo, locator)

		currentTree := getTestTrees()[0]
		treeCluster := getTestTreeClusters()[0]
		currentTree.TreeCluster = treeCluster
		sensor := getTestSensors()[0]
		currentTree.Sensor = sensor

		// Mock expectations
		treeRepo.EXPECT().GetByID(ctx, id).Return(currentTree, nil)
		treeClusterRepo.EXPECT().GetByID(ctx, currentTree.TreeCluster.ID).Return(treeCluster, nil)
		sensorRepo.EXPECT().GetByID(ctx, currentTree.Sensor.ID).Return(sensor, nil)
		treeRepo.EXPECT().Update(ctx, id,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything).Return(updatedTree, nil)
		locator.EXPECT().UpdateCluster(ctx, &currentTree.TreeCluster.ID).Return(nil)
		locator.EXPECT().UpdateCluster(ctx, &updatedTree.TreeCluster.ID).Return(nil)

		// When
		result, err := svc.Update(ctx, id, treeUpdate)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, updatedTree, result)
	})

	t.Run("should return validation error", func(t *testing.T) {
		// Given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		imageRepo := storageMock.NewMockImageRepository(t)
		treeClusterRepo := storageMock.NewMockTreeClusterRepository(t)
		locator := service.NewMockGeoClusterLocator(t)

		svc := NewTreeService(treeRepo, sensorRepo, imageRepo, treeClusterRepo, locator)

		invalidTreeUpdate := &entities.TreeUpdate{
			Latitude:     0,
			Longitude:    0,
			PlantingYear: -2013,
		}

		// When
		result, err := svc.Update(ctx, 1, invalidTreeUpdate)

		// Then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validation error")
	})
	t.Run("should return error if tree not found", func(t *testing.T) {
		// Given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		imageRepo := storageMock.NewMockImageRepository(t)
		treeClusterRepo := storageMock.NewMockTreeClusterRepository(t)
		locator := service.NewMockGeoClusterLocator(t)

		svc := NewTreeService(treeRepo, sensorRepo, imageRepo, treeClusterRepo, locator)

		expectedError := storage.ErrTreeNotFound

		// Mock expectations
		treeRepo.EXPECT().GetByID(ctx, id).Return(nil, expectedError)

		// When
		result, err := svc.Update(ctx, id, treeUpdate)

		// Then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, handleError(expectedError).Error())
	})

	t.Run("should return error if TreeCluster not found", func(t *testing.T) {
		// Given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		imageRepo := storageMock.NewMockImageRepository(t)
		treeClusterRepo := storageMock.NewMockTreeClusterRepository(t)
		locator := service.NewMockGeoClusterLocator(t)

		svc := NewTreeService(treeRepo, sensorRepo, imageRepo, treeClusterRepo, locator)

		expectedError := storage.ErrTreeClusterNotFound

		currentTree := getTestTrees()[0]
		treeCluster := getTestTreeClusters()[0]
		currentTree.TreeCluster = treeCluster

		// Mock expectations
		treeRepo.EXPECT().GetByID(ctx, id).Return(currentTree, nil)

		treeClusterRepo.EXPECT().GetByID(ctx, currentTree.TreeCluster.ID).Return(nil, expectedError)

		// When
		result, err := svc.Update(ctx, id, treeUpdate)

		// Then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, expectedError.Error())
	})

	t.Run("should return error if Sensor not found", func(t *testing.T) {
		// Given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		imageRepo := storageMock.NewMockImageRepository(t)
		treeClusterRepo := storageMock.NewMockTreeClusterRepository(t)
		locator := service.NewMockGeoClusterLocator(t)

		svc := NewTreeService(treeRepo, sensorRepo, imageRepo, treeClusterRepo, locator)

		expectedError := storage.ErrSensorNotFound

		currentTree := getTestTrees()[0]
		sensor := getTestSensors()[0]
		currentTree.Sensor = sensor
		treeCluster := getTestTreeClusters()[0]
		currentTree.TreeCluster = treeCluster

		// Mock expectations
		treeRepo.EXPECT().GetByID(ctx, id).Return(currentTree, nil)
		treeClusterRepo.EXPECT().GetByID(ctx, currentTree.TreeCluster.ID).Return(treeCluster, nil)
		sensorRepo.EXPECT().GetByID(ctx, currentTree.Sensor.ID).Return(nil, expectedError)

		// When
		result, err := svc.Update(ctx, id, treeUpdate)

		// Then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, expectedError.Error())
	})
	t.Run("should return error if updating tree fails", func(t *testing.T) {
		// Given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		imageRepo := storageMock.NewMockImageRepository(t)
		treeClusterRepo := storageMock.NewMockTreeClusterRepository(t)
		locator := service.NewMockGeoClusterLocator(t)

		svc := NewTreeService(treeRepo, sensorRepo, imageRepo, treeClusterRepo, locator)

		expectedError := errors.New("update failed")

		currentTree := getTestTrees()[0]
		treeCluster := getTestTreeClusters()[0]
		currentTree.TreeCluster = treeCluster
		sensor := getTestSensors()[0]
		currentTree.Sensor = sensor

		// Mock expectations
		treeRepo.EXPECT().GetByID(ctx, id).Return(currentTree, nil)
		treeClusterRepo.EXPECT().GetByID(ctx, currentTree.TreeCluster.ID).Return(treeCluster, nil)
		sensorRepo.EXPECT().GetByID(ctx, currentTree.Sensor.ID).Return(sensor, nil)
		treeRepo.EXPECT().Update(ctx, id,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything).Return(nil, expectedError)

		// When
		result, err := svc.Update(ctx, id, treeUpdate)

		// Then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, handleError(expectedError).Error())
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

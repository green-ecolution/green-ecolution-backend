package sensor

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	treeUtils "github.com/green-ecolution/green-ecolution-backend/internal/service/domain/tree"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	storageMock "github.com/green-ecolution/green-ecolution-backend/internal/storage/_mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestSensorService_GetAll(t *testing.T) {
	t.Run("should return all sensor", func(t *testing.T) {
		// given
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := NewSensorService(sensorRepo, treeRepo, flowerbedRepo)

		// when
		sensorRepo.EXPECT().GetAll(context.Background()).Return(TestSensorList, nil)
		sensors, err := svc.GetAll(context.Background())

		// then
		assert.NoError(t, err)
		assert.Equal(t, TestSensorList, sensors)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// given
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := NewSensorService(sensorRepo, treeRepo, flowerbedRepo)

		sensorRepo.EXPECT().GetAll(context.Background()).Return(nil, storage.ErrSensorNotFound)
		sensors, err := svc.GetAll(context.Background())

		// then
		assert.Error(t, err)
		assert.Nil(t, sensors)
		assert.EqualError(t, err, "500: sensor not found")
	})
}

func TestSensorService_GetByID(t *testing.T) {
	t.Run("should return sensor when found", func(t *testing.T) {
		// given
		id := "sensor-1"
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := NewSensorService(sensorRepo, treeRepo, flowerbedRepo)

		sensorRepo.EXPECT().GetByID(context.Background(), id).Return(TestSensor, nil)

		// when
		sensor, err := svc.GetByID(context.Background(), id)

		// then
		assert.NoError(t, err)
		assert.Equal(t, TestSensor, sensor)
	})

	t.Run("should return error if sensor not found", func(t *testing.T) {
		// given
		id := "sensor-1"
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := NewSensorService(sensorRepo, treeRepo, flowerbedRepo)

		expectedErr := storage.ErrEntityNotFound
		sensorRepo.EXPECT().GetByID(context.Background(), id).Return(nil, expectedErr)

		// when
		sensor, err := svc.GetByID(context.Background(), id)

		// then
		assert.Error(t, err)
		assert.Nil(t, sensor)
		assert.EqualError(t, err, "404: sensor not found")
	})
}

func TestSensorService_Create(t *testing.T) {
	newSensor := &entities.SensorCreate{
		ID:        "sensor-1",
		Status:    entities.SensorStatusOnline,
		Data:      TestSensor.Data,
		Latitude:  9.446741,
		Longitude: 54.801539,
	}

	t.Run("should successfully create a new sensor", func(t *testing.T) {
		// given
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := NewSensorService(sensorRepo, treeRepo, flowerbedRepo)

		sensorRepo.EXPECT().Create(
			context.Background(),
			mock.Anything,
			mock.Anything,
		).Return(TestSensor, nil)

		// when
		result, err := svc.Create(context.Background(), newSensor)

		// then
		assert.NoError(t, err)
		assert.Equal(t, TestSensor, result)
	})

	t.Run("should successfully create a new sensor without data", func(t *testing.T) {
		// given
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := NewSensorService(sensorRepo, treeRepo, flowerbedRepo)

		newSensor.Data = []*entities.SensorData{}

		sensorRepo.EXPECT().Create(
			context.Background(),
			mock.Anything,
			mock.Anything,
		).Return(TestSensor, nil)

		// when
		result, err := svc.Create(context.Background(), newSensor)

		// then
		assert.NoError(t, err)
		assert.Equal(t, TestSensor, result)
	})

	t.Run("should return validation error on no status", func(t *testing.T) {
		// given
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := NewSensorService(sensorRepo, treeRepo, flowerbedRepo)

		newSensor.Status = ""

		// when
		result, err := svc.Create(context.Background(), newSensor)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validation error")
	})

	t.Run("should return validation error on invalid sensor id", func(t *testing.T) {
		// given
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := NewSensorService(sensorRepo, treeRepo, flowerbedRepo)

		newSensor.Status = entities.SensorStatusOffline
		newSensor.ID = ""

		// when
		result, err := svc.Create(context.Background(), newSensor)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validation error")
	})

	t.Run("should return validation error on invalid latitude and longitude", func(t *testing.T) {
		// given
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := NewSensorService(sensorRepo, treeRepo, flowerbedRepo)

		newSensor.ID = "sensor-23"
		newSensor.Status = entities.SensorStatusOffline
		newSensor.Latitude = 200
		newSensor.Longitude = 100

		// when
		result, err := svc.Create(context.Background(), newSensor)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validation error")
	})

	t.Run("should return an error when creating sensor fails", func(t *testing.T) {
		// given
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := NewSensorService(sensorRepo, treeRepo, flowerbedRepo)
		expectedErr := errors.New("Failed to create sensor")

		newSensor.ID = "sensor-23"
		newSensor.Status = entities.SensorStatusOffline
		newSensor.Latitude = 9.446741
		newSensor.Longitude = 54.801539

		sensorRepo.EXPECT().Create(
			context.Background(),
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedErr)

		// when
		result, err := svc.Create(context.Background(), newSensor)

		// then
		assert.Nil(t, result)
		assert.EqualError(t, err, "500: Failed to create sensor")
	})
}

func TestSensorService_Update(t *testing.T) {
	updateSensor := &entities.SensorUpdate{
		Status:    entities.SensorStatusOnline,
		Data:      TestSensor.Data,
		Latitude:  9.446741,
		Longitude: 54.801539,
	}

	t.Run("should successfully update a sensor", func(t *testing.T) {
		// given
		id := "sensor-1"
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := NewSensorService(sensorRepo, treeRepo, flowerbedRepo)

		sensorRepo.EXPECT().GetByID(context.Background(), id).Return(TestSensor, nil)

		sensorRepo.EXPECT().Update(
			context.Background(),
			id,
			mock.Anything,
			mock.Anything,
		).Return(TestSensor, nil)

		// when
		result, err := svc.Update(context.Background(), id, updateSensor)

		// then
		assert.NoError(t, err)
		assert.Equal(t, TestSensor, result)
	})

	t.Run("should return an error when sensor ID does not exist", func(t *testing.T) {
		// given
		id := "notFoundID"
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := NewSensorService(sensorRepo, treeRepo, flowerbedRepo)
		expectedErr := errors.New("failed to update cluster")

		sensorRepo.EXPECT().GetByID(context.Background(), id).Return(nil, expectedErr)

		// when
		result, err := svc.Update(context.Background(), id, updateSensor)

		// then
		assert.Nil(t, result)
		assert.EqualError(t, err, "500: failed to update cluster")
	})

	t.Run("should return an error when the update fails", func(t *testing.T) {
		// given
		id := "sensor-1"
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := NewSensorService(sensorRepo, treeRepo, flowerbedRepo)
		expectedErr := errors.New("failed to update cluster")

		sensorRepo.EXPECT().GetByID(context.Background(), id).Return(TestSensor, nil)

		sensorRepo.EXPECT().Update(
			context.Background(),
			id,
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedErr)

		// when
		result, err := svc.Update(context.Background(), id, updateSensor)

		// then
		assert.Nil(t, result)
		assert.EqualError(t, err, "500: failed to update cluster")
	})

	t.Run("should return validation error on invalid latitude and longitude", func(t *testing.T) {
		// given
		id := "sensor-1"
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := NewSensorService(sensorRepo, treeRepo, flowerbedRepo)

		updateSensor.Latitude = 200
		updateSensor.Longitude = 200

		// when
		result, err := svc.Update(context.Background(), id, updateSensor)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validation error")
	})
}

func TestSensorService_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully delete a sensor", func(t *testing.T) {
		// given
		id := "sensor-1"
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := NewSensorService(sensorRepo, treeRepo, flowerbedRepo)

		sensorRepo.EXPECT().GetByID(ctx, id).Return(TestSensor, nil)
		treeRepo.EXPECT().UnlinkSensorID(ctx, id).Return(nil)
		flowerbedRepo.EXPECT().UnlinkSensorID(ctx, id).Return(nil)
		sensorRepo.EXPECT().Delete(ctx, id).Return(nil)

		// when
		err := svc.Delete(ctx, id)

		// then
		assert.NoError(t, err)
	})

	t.Run("should return error if sensor not found", func(t *testing.T) {
		// given
		id := "sensor-1"
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := NewSensorService(sensorRepo, treeRepo, flowerbedRepo)

		expectedErr := storage.ErrEntityNotFound
		sensorRepo.EXPECT().GetByID(ctx, id).Return(nil, expectedErr)

		// when
		err := svc.Delete(ctx, id)

		// then
		assert.Error(t, err)
		assert.EqualError(t, err, "404: sensor not found")
	})

	t.Run("should return error if unlinking sensor ID on tree fails", func(t *testing.T) {
		// given
		id := "sensor-1"
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := NewSensorService(sensorRepo, treeRepo, flowerbedRepo)

		expectedErr := errors.New("failed to unlink")

		sensorRepo.EXPECT().GetByID(ctx, id).Return(TestSensor, nil)
		treeRepo.EXPECT().UnlinkSensorID(ctx, id).Return(expectedErr)

		// when
		err := svc.Delete(ctx, id)

		// then
		assert.Error(t, err)
		assert.EqualError(t, err, "500: failed to unlink")
	})

	t.Run("should return error if unlinking sensor ID on flowerbed fails", func(t *testing.T) {
		// given
		id := "sensor-1"
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := NewSensorService(sensorRepo, treeRepo, flowerbedRepo)
		expectedErr := errors.New("failed to unlink")

		sensorRepo.EXPECT().GetByID(ctx, id).Return(TestSensor, nil)
		treeRepo.EXPECT().UnlinkSensorID(ctx, id).Return(nil)
		flowerbedRepo.EXPECT().UnlinkSensorID(ctx, id).Return(expectedErr)

		// when
		err := svc.Delete(ctx, id)

		// then
		assert.Error(t, err)
		assert.EqualError(t, err, "500: failed to unlink")
	})

	t.Run("should return error if deleting sensor fails", func(t *testing.T) {
		// given
		id := "sensor-1"
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := NewSensorService(sensorRepo, treeRepo, flowerbedRepo)
		expectedErr := errors.New("failed to delete")

		sensorRepo.EXPECT().GetByID(ctx, id).Return(TestSensor, nil)
		treeRepo.EXPECT().UnlinkSensorID(ctx, id).Return(nil)
		flowerbedRepo.EXPECT().UnlinkSensorID(ctx, id).Return(nil)
		sensorRepo.EXPECT().Delete(ctx, id).Return(expectedErr)

		// when
		err := svc.Delete(ctx, id)

		// then
		assert.Error(t, err)
		assert.EqualError(t, err, "500: failed to delete")
	})
}

func TestSensorService_MapSensorToTree(t *testing.T) {
	t.Run("should successfully map sensor to the nearest tree", func(t *testing.T) {
		// given
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)

		svc := SensorService{sensorRepo: sensorRepo, treeRepo: treeRepo, validator: validator.New()}
		testSensor := TestSensorNearestTree
		testTree := treeUtils.TestNearestTree

		treeRepo.EXPECT().
			FindNearestTree(context.Background(), mock.Anything, mock.Anything).
			Return(testTree, nil)

		treeRepo.EXPECT().
			Update(context.Background(), testTree.ID, mock.Anything).
			Return(testTree, nil)

		// when
		err := svc.MapSensorToTree(context.Background(), testSensor)

		// then
		assert.NoError(t, err)
	})

	t.Run("should return error if sensor is nil", func(t *testing.T) {
		// given
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		svc := SensorService{sensorRepo: sensorRepo, treeRepo: treeRepo, validator: validator.New()}

		// when
		err := svc.MapSensorToTree(context.Background(), nil)

		// then
		assert.Error(t, err)
		assert.EqualError(t, err, "sensor cannot be nil")
	})

	t.Run("should return error if nearest tree is not found", func(t *testing.T) {
		// given
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		svc := SensorService{sensorRepo: sensorRepo, treeRepo: treeRepo, validator: validator.New()}
		testSensor := TestSensorNearestTree

		treeRepo.EXPECT().
			FindNearestTree(context.Background(), testSensor.Latitude, testSensor.Longitude).
			Return(nil, errors.New("tree not found"))

		// when
		err := svc.MapSensorToTree(context.Background(), testSensor)

		// then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tree not found")
	})

	t.Run("should return error if updating tree with sensor fails", func(t *testing.T) {
		// given
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		svc := SensorService{sensorRepo: sensorRepo, treeRepo: treeRepo, validator: validator.New()}
		testSensor := TestSensorNearestTree
		testTree := treeUtils.TestNearestTree

		treeRepo.EXPECT().
			FindNearestTree(context.Background(), mock.Anything, mock.Anything).
			Return(testTree, nil)

		treeRepo.EXPECT().
			Update(context.Background(), testTree.ID, mock.Anything).
			Return(nil, errors.New("update failed"))

		// when
		err := svc.MapSensorToTree(context.Background(), testSensor)

		// then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update failed")
	})
}

func TestReady(t *testing.T) {
	t.Run("should return true if the service is ready", func(t *testing.T) {
		// given
		repo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := NewSensorService(repo, treeRepo, flowerbedRepo)

		// when
		ready := svc.Ready()

		// then
		assert.True(t, ready)
	})

	t.Run("should return false if the service is not ready", func(t *testing.T) {
		// give
		svc := NewSensorService(nil, nil, nil)

		// when
		ready := svc.Ready()

		// then
		assert.False(t, ready)
	})
}

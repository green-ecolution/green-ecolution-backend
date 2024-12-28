package sensor_test

import (
	"context"
	"errors"
	"testing"

	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/service/domain/sensor"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/stretchr/testify/mock"

	storageMock "github.com/green-ecolution/green-ecolution-backend/internal/storage/_mock"
	"github.com/stretchr/testify/assert"
)

func TestNewSensorService(t *testing.T) {
	t.Run("should create a new service", func(t *testing.T) {
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := sensor.NewSensorService(sensorRepo, treeRepo, flowerbedRepo)
		assert.NotNil(t, svc)
	})
}

func TestSensorService_HandleMessage(t *testing.T) {
	t.Run("should process MQTT payload successfully", func(t *testing.T) {
		// given
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := sensor.NewSensorService(sensorRepo, treeRepo, flowerbedRepo)

		testPayLoad := TestListMQTTPayload[0]
		insertData := &domain.SensorData{
			Data: testPayLoad,
		}

		sensorRepo.EXPECT().GetByID(context.Background(), testPayLoad.DeviceID).Return(TestSensor, nil)
		sensorRepo.EXPECT().Update(context.Background(),
			TestSensor.ID,
			mock.Anything,
			mock.Anything,
			mock.Anything).Return(TestSensor, nil)
		sensorRepo.EXPECT().InsertSensorData(context.Background(), insertData, testPayLoad.DeviceID).Return(nil)
		sensorRepo.EXPECT().GetLastSensorDataByID(context.Background(), TestSensor.ID).Return(TestSensorData[0], nil)

		// when
		sensorData, err := svc.HandleMessage(context.Background(), testPayLoad)
		sensor, errGetSens := sensorRepo.GetByID(context.Background(), TestSensor.ID)

		// then
		assert.NoError(t, err)
		assert.NoError(t, errGetSens)
		assert.NotNil(t, sensorData)
		assert.NotEmpty(t, sensorData)
		assert.Equal(t, sensorData.Data, insertData.Data)
		assert.NotNil(t, sensor)
		assert.Equal(t, sensor.Latitude, TestSensor.Latitude)
		assert.Equal(t, sensor.Longitude, TestSensor.Longitude)
		assert.Equal(t, sensor.Status, domain.SensorStatusOnline)
	})

	t.Run("should return error if sensor update fails", func(t *testing.T) {
		// given
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := sensor.NewSensorService(sensorRepo, treeRepo, flowerbedRepo)

		testPayload := TestListMQTTPayload[0]

		sensorRepo.EXPECT().GetByID(context.Background(), testPayload.DeviceID).Return(TestSensor, nil)
		sensorRepo.EXPECT().Update(context.Background(),
			TestSensor.ID,
			mock.Anything,
			mock.Anything,
			mock.Anything).
			Return(nil, errors.New("update error"))

		// when
		sensorData, err := svc.HandleMessage(context.Background(), testPayload)

		// then
		assert.Error(t, err)
		assert.Nil(t, sensorData)
		assert.Contains(t, err.Error(), "update error")
	})

	t.Run("should process MQTT payload and create a new sensor if not found", func(t *testing.T) {
		// given
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := sensor.NewSensorService(sensorRepo, treeRepo, flowerbedRepo)

		testPayLoad := TestListMQTTPayload[0]
		insertData := &domain.SensorData{
			Data: testPayLoad,
		}
	
		sensorRepo.EXPECT().GetByID(context.Background(), testPayLoad.DeviceID).Return(nil, storage.ErrSensorNotFound).Once()
		sensorRepo.EXPECT().Create(context.Background(),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything).
			Return(TestSensor, nil).Once()
		sensorRepo.EXPECT().InsertSensorData(context.Background(), insertData, TestSensor.ID).Return(nil).Once()
		sensorRepo.EXPECT().GetLastSensorDataByID(context.Background(), TestSensor.ID).Return(TestSensorData[0], nil).Once()
		sensorRepo.EXPECT().GetByID(context.Background(), TestSensor.ID).Return(TestSensor, nil).Once()
	
		// when
		sensorData, err := svc.HandleMessage(context.Background(), testPayLoad)
		sensor, errCreateSens := sensorRepo.GetByID(context.Background(), TestSensor.ID)
	
		// then
		assert.NoError(t, err)
		assert.NoError(t, errCreateSens)
		assert.NotNil(t, sensorData)
		assert.NotEmpty(t, sensorData)
		assert.Equal(t, sensorData.Data, insertData.Data)
		assert.NotNil(t, sensor)
		assert.Equal(t, sensor.Latitude, TestSensor.Latitude)
		assert.Equal(t, sensor.Longitude, TestSensor.Longitude)
		assert.Equal(t, sensor.Status, domain.SensorStatusOnline)
	})
	
	t.Run("should return error if sensor creation fails", func(t *testing.T) {
		// given
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := sensor.NewSensorService(sensorRepo, treeRepo, flowerbedRepo)

		testPayload := TestListMQTTPayload[0]

		sensorRepo.EXPECT().GetByID(context.Background(), testPayload.DeviceID).Return(nil, storage.ErrSensorNotFound)
		sensorRepo.EXPECT().Create(context.Background(),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything).
			Return(nil, errors.New("create error"))

		// when
		sensorData, err := svc.HandleMessage(context.Background(), testPayload)

		// then
		assert.Error(t, err)
		assert.Nil(t, sensorData)
		assert.Contains(t, err.Error(), "create error")
	})

	t.Run("should return error when payload is nil", func(t *testing.T) {
		// given
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := sensor.NewSensorService(sensorRepo, treeRepo, flowerbedRepo)

		// when
		result, err := svc.HandleMessage(context.Background(), nil)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("should return validation error for invalid latitude", func(t *testing.T) {
		// given
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := sensor.NewSensorService(sensorRepo, treeRepo, flowerbedRepo)

		// when
		result, err := svc.HandleMessage(context.Background(), TestMQTTPayLoadInvalidLat)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), service.ErrValidation.Error())
	})

	t.Run("should return validation error for invalid longitude", func(t *testing.T) {
		// given
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := sensor.NewSensorService(sensorRepo, treeRepo, flowerbedRepo)

		// when
		result, err := svc.HandleMessage(context.Background(), TestMQTTPayLoadInvalidLong)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), service.ErrValidation.Error())
	})

	t.Run("should return error if InsertSensorData fails", func(t *testing.T) {
		// given
		sensorRepo := storageMock.NewMockSensorRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		flowerbedRepo := storageMock.NewMockFlowerbedRepository(t)
		svc := sensor.NewSensorService(sensorRepo, treeRepo, flowerbedRepo)

		testPayLoad := TestListMQTTPayload[0]
		insertData := &domain.SensorData{
			Data: testPayLoad,
		}

		sensorRepo.EXPECT().GetByID(context.Background(), testPayLoad.DeviceID).Return(TestSensor, nil)
		sensorRepo.EXPECT().Update(context.Background(),
			TestSensor.ID,
			mock.Anything,
			mock.Anything,
			mock.Anything).Return(TestSensor, nil)
			sensorRepo.EXPECT().InsertSensorData(context.Background(), insertData, testPayLoad.DeviceID).Return(errors.New("insert error"))

		// when
		sensorData, err := svc.HandleMessage(context.Background(), testPayLoad)

		// then
		assert.Error(t, err)
		assert.Nil(t, sensorData)
		assert.Contains(t, err.Error(), "insert error")
	})
}

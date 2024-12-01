package sensor

import (
	"context"
	"errors"
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/stretchr/testify/mock"
	"testing"

	storageMock "github.com/green-ecolution/green-ecolution-backend/internal/storage/_mock"
	"github.com/stretchr/testify/assert"
)

func TestNewSensorService(t *testing.T) {
	repo := storageMock.NewMockSensorRepository(t)
	t.Run("should create a new service", func(t *testing.T) {
		svc := NewMqttService(repo)
		assert.NotNil(t, svc)
	})
}

func TestSensorService_HandleMessage(t *testing.T) {
	t.Run("should process MQTT payload successfully", func(t *testing.T) {
		// given
		sensorRepo := storageMock.NewMockSensorRepository(t)
		svc := NewMqttService(sensorRepo)

		testPayLoad := TestListMQTTPayload[0]
		insertData := []*domain.SensorData{
			{
				Data: testPayLoad,
			},
		}

		sensorRepo.EXPECT().GetByID(context.Background(), testPayLoad.DeviceID).Return(TestSensor, nil)
		sensorRepo.EXPECT().InsertSensorData(context.Background(), insertData, testPayLoad.DeviceID).Return(insertData, nil)
		sensorRepo.EXPECT().Update(context.Background(),
			TestSensor.ID,
			mock.Anything,
			mock.Anything).Return(TestSensor, nil)
		sensorRepo.EXPECT().GetSensorDataByID(context.Background(), TestSensor.ID).Return(TestSensorData, nil)

		// when
		sensorData, err := svc.HandleMessage(context.Background(), testPayLoad)
		sensor, errGetSens := svc.sensorRepo.GetByID(context.Background(), TestSensor.ID)

		// then
		assert.NoError(t, err)
		assert.NoError(t, errGetSens)
		assert.NotNil(t, sensorData)
		assert.NotEmpty(t, sensorData)
		assert.Equal(t, sensorData[0].Data, insertData[0].Data)
		assert.NotNil(t, sensor)
		assert.Equal(t, sensor.Latitude, TestSensor.Latitude)
		assert.Equal(t, sensor.Longitude, TestSensor.Longitude)
	})

	t.Run("should return error when payload is nil", func(t *testing.T) {
		// given
		sensorRepo := storageMock.NewMockSensorRepository(t)
		svc := NewMqttService(sensorRepo)

		// when
		result, err := svc.HandleMessage(context.Background(), nil)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("should return validation error for invalid latitude", func(t *testing.T) {
		// given
		sensorRepo := storageMock.NewMockSensorRepository(t)
		svc := NewMqttService(sensorRepo)

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
		svc := NewMqttService(sensorRepo)

		// when
		result, err := svc.HandleMessage(context.Background(), TestMQTTPayLoadInvalidLong)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), service.ErrValidation.Error())
	})

	t.Run("should return error if sensor is not found", func(t *testing.T) {
		// given
		sensorRepo := storageMock.NewMockSensorRepository(t)
		svc := NewMqttService(sensorRepo)

		testPayLoad := TestListMQTTPayload[0]
		testPayLoad.DeviceID = "notFoundSensorID"

		sensorRepo.EXPECT().GetByID(context.Background(), testPayLoad.DeviceID).Return(nil, storage.ErrSensorNotFound)

		// when
		sensorData, err := svc.HandleMessage(context.Background(), testPayLoad)
		sensor, errGetSens := svc.sensorRepo.GetByID(context.Background(), testPayLoad.DeviceID)

		// then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), storage.ErrSensorNotFound.Error())
		assert.Nil(t, sensorData)
		assert.Error(t, errGetSens)
		assert.Nil(t, sensor)
	})

	t.Run("should return error if InsertSensorData fails", func(t *testing.T) {
		// given
		sensorRepo := storageMock.NewMockSensorRepository(t)
		svc := NewMqttService(sensorRepo)

		testPayLoad := TestListMQTTPayload[0]
		insertData := []*domain.SensorData{
			{
				Data: testPayLoad,
			},
		}

		sensorRepo.EXPECT().GetByID(context.Background(), testPayLoad.DeviceID).Return(TestSensor, nil)
		sensorRepo.EXPECT().InsertSensorData(context.Background(), insertData, testPayLoad.DeviceID).Return(insertData, errors.New("insert error"))

		// when
		sensorData, err := svc.HandleMessage(context.Background(), testPayLoad)

		// then
		assert.Error(t, err)
		assert.Nil(t, sensorData)
		assert.Contains(t, err.Error(), "insert error")
	})
}

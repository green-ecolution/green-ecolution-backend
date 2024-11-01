package sensor

import (
	"context"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sensorUtils "github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/sensor"
	"github.com/stretchr/testify/assert"
)

func TestSensorRepository_Create(t *testing.T) {
	suite.ResetDB(t)

	t.Run("should create sensor", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.Create(
			context.Background(), 
			WithStatus(sensorUtils.TestSensor.Status),
			WithData(sensorUtils.TestSensor.Data),
		)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, sensorUtils.TestSensor.Status, got.Status)
		assert.NotZero(t, got.ID)
	})

	t.Run("should create sensor with empty data and unknown status", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.Create(context.Background())

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, entities.SensorStatusUnknown, got.Status)
		assert.Equal(t, []*entities.SensorData(nil), got.Data)
		assert.NotZero(t, got.ID)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.Create(ctx)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestSensorRepository_InsertSensorData(t *testing.T) {
	suite.ResetDB(t)

	t.Run("should insert sensor data successfully", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())
		sensor := &entities.Sensor{
			Status: entities.SensorStatusOnline,
		}
		
		_, err := r.Create(context.Background(), WithStatus(sensor.Status))
		assert.NoError(t, err)

		data := []*entities.SensorData{
			{
				ID:   1,
				Data: sensorUtils.TestMqttPayload,
			},
		}

		// when
		got, err := r.InsertSensorData(context.Background(), data)

		// then
		assert.NoError(t, err)
		assert.Equal(t, data, got)
	})

	t.Run("should return error when data is empty", func(t *testing.T) {
		r := NewSensorRepository(suite.Store, defaultSensorMappers())
		sensor := &entities.Sensor{
			Status: entities.SensorStatusOnline,
		}
		
		_, err := r.Create(context.Background(), WithStatus(sensor.Status))
		assert.NoError(t, err)

		data := []*entities.SensorData{}

		// when
		got, err := r.InsertSensorData(context.Background(), data)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when data is nil", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())
	
		// when
		got, err := r.InsertSensorData(context.Background(), nil)
	
		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

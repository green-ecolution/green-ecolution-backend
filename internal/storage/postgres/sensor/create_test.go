package sensor

import (
	"context"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/storage"

	"github.com/green-ecolution/green-ecolution-backend/internal/utils"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sensorUtils "github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/sensor"
	"github.com/stretchr/testify/assert"
)

func TestSensorRepository_Create(t *testing.T) {
	t.Run("should create sensor", func(t *testing.T) {
		suite.ResetDB(t)
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())
		sensorID := "sensor-4"
		sensorData := sensorUtils.TestSensorData[1]
		sensorData.SensorID = utils.P(sensorID)
		sensorDataList := []*entities.SensorData{sensorData}

		// when
		got, err := r.Create(
			context.Background(),
			WithSensorID(sensorID),
			WithStatus(entities.SensorStatusOnline),
			WithData(sensorDataList),
		)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, float64(0), got.Latitude)
		assert.Equal(t, float64(0), got.Longitude)
		assert.Equal(t, entities.SensorStatusOnline, got.Status)
		assert.NotZero(t, got.ID)
	})

	t.Run("should create sensor with empty data and unknown status", func(t *testing.T) {
		suite.ResetDB(t)
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.Create(context.Background(),
			WithSensorID("sensor-5"),
			WithLatitude(54.801539),
			WithLongitude(9.446741))

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, entities.SensorStatusUnknown, got.Status)
		assert.Equal(t, []*entities.SensorData(nil), got.Data)
		assert.Equal(t, 54.801539, got.Latitude)
		assert.Equal(t, 9.446741, got.Longitude)
		assert.NotZero(t, got.ID)
	})

	t.Run("should return error if latitude is out of bounds", func(t *testing.T) {
		suite.ResetDB(t)
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.Create(context.Background(),
			WithSensorID("sensor-5"),
			WithLatitude(-200),
			WithLongitude(0))

		// then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), storage.ErrInvalidLatitude.Error())
		assert.Nil(t, got)
	})

	t.Run("should return error if longitude is out of bounds", func(t *testing.T) {
		suite.ResetDB(t)
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.Create(context.Background(),
			WithSensorID("sensor-5"),
			WithLatitude(0),
			WithLongitude(200))

		// then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), storage.ErrInvalidLongitude.Error())
		assert.Nil(t, got)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.Create(ctx, WithSensorID("sensor-5"))

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestSensorRepository_InsertSensorData(t *testing.T) {
	t.Run("should insert sensor data successfully", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())
		sensor := &entities.Sensor{
			Status: entities.SensorStatusOnline,
		}

		_, err := r.Create(context.Background(), WithSensorID(sensorUtils.TestSensorID), WithStatus(sensor.Status))
		assert.NoError(t, err)

		data := []*entities.SensorData{
			{
				ID:       1,
				SensorID: utils.P(sensorUtils.TestSensorID),
				Data:     sensorUtils.TestMqttPayload,
			},
		}

		// when
		got, err := r.InsertSensorData(context.Background(), data)

		// then
		assert.NoError(t, err)
		assert.Equal(t, data, got)
	})

	t.Run("should return error when data is empty", func(t *testing.T) {
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())
		sensor := &entities.Sensor{
			Status: entities.SensorStatusOnline,
		}

		_, err := r.Create(context.Background(), WithSensorID(sensorUtils.TestSensorID), WithStatus(sensor.Status))
		assert.NoError(t, err)

		var data []*entities.SensorData

		// when
		got, err := r.InsertSensorData(context.Background(), data)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when data is nil", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.InsertSensorData(context.Background(), nil)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

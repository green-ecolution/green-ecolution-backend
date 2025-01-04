package sensor

import (
	"context"
	"testing"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/storage"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
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
			WithSensorID(input.ID),
			WithLatitude(input.Latitude),
			WithLongitude(input.Longitude),
			WithStatus(input.Status),
			WithLatestData(input.LatestData),
		)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, input.Latitude, got.Latitude)
		assert.Equal(t, input.Longitude, got.Longitude)
		assert.Equal(t, input.Status, got.Status)
		assert.NotZero(t, got.ID)

		// assert latest data
		assert.NotZero(t, got.LatestData.UpdatedAt)
		assert.NotZero(t, got.LatestData.CreatedAt)
		assert.Equal(t, input.LatestData.Data, got.LatestData.Data)
	})
	t.Run("should handle error gracefully when creating a sensor with empty data and unknown status", func(t *testing.T) {
		// Arrange
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// Valid latitude and longitude
		latitude := 0.0
		longitude := 0.0

		// Act
		got, err := r.Create(
			context.Background(),
			WithSensorID("sensor-126"),
			WithLatitude(latitude),
			WithLongitude(longitude),
		)

		// Handle error with context
		if err != nil {
			err = suite.Store.HandleError(err, "Error occurred while creating a sensor with empty data")
		}

		// Assert: error is classified
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.EqualError(t, err, "UnexpectedError Error occurred while creating a sensor with empty data (at internal/storage/postgres/sensor/create_test.go:62)")
	})

	t.Run("return errors properly when creating a sensor with invalid parameters", func(t *testing.T) {
		// Arrange
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// Invalid latitude and longitude
		invalidLatitude := -200.0
		invalidLongitude := 200.0

		// Act
		got, err := r.Create(
			context.Background(),
			WithSensorID("sensor-127"),
			WithLatitude(invalidLatitude),
			WithLongitude(invalidLongitude),
		)

		// Handle error with context
		if err != nil {
			err = suite.Store.HandleError(err, "Invalid sensor parameters provided")
		}

		// Assert: error is classified
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.EqualError(t, err, "UnexpectedError Invalid sensor parameters provided (at internal/storage/postgres/sensor/create_test.go:89)")
	})

	t.Run("should return error if latitude is out of bounds", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.Create(context.Background(),
			WithSensorID("sensor-125"),
			WithLatitude(-200),
			WithLongitude(input.Longitude))

		// then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), storage.ErrInvalidLatitude.Error())
		assert.Nil(t, got)
	})

	t.Run("should return error if longitude is out of bounds", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.Create(context.Background(),
			WithSensorID("sensor-125"),
			WithLatitude(input.Latitude),
			WithLongitude(200))

		// then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), storage.ErrInvalidLongitude.Error())
		assert.Nil(t, got)
	})

	t.Run("should return error if sensor id is invalid", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.Create(context.Background(),
			WithLatitude(input.Latitude),
			WithLongitude(input.Longitude))

		// then
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "sensor id cannot be empty")
		assert.Nil(t, got)
	})

	t.Run("should return error if sensor with same id already exists", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.Create(
			context.Background(),
			WithSensorID(input.ID),
			WithLatitude(input.Latitude),
			WithLongitude(input.Longitude),
			WithStatus(input.Status),
			WithLatestData(input.LatestData),
		)

		// then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sensor with same ID already exists")
		assert.Nil(t, got)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
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
	suite.ResetDB(t)

	t.Run("should insert sensor data successfully", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		_, err := r.Create(context.Background(),
			WithSensorID(input.ID),
			WithLatitude(input.Latitude),
			WithLongitude(input.Longitude),
			WithStatus(input.Status))

		assert.Error(t, err)

		// when
		err = r.InsertSensorData(context.Background(), input.LatestData, input.ID)

		// then
		assert.NoError(t, err)
	})

	t.Run("should return error when data is empty", func(t *testing.T) {
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		_, err := r.Create(context.Background(),
			WithSensorID("sensor-124"),
			WithLatitude(input.Latitude),
			WithLongitude(input.Longitude),
			WithStatus(input.Status))

		assert.Error(t, err)

		// when
		err = r.InsertSensorData(context.Background(), &entities.SensorData{}, input.ID)

		// then
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "latest data cannot be empty")
	})

	t.Run("should return error when data is nil", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		err := r.InsertSensorData(context.Background(), nil, "sensor-1")

		// then
		assert.Error(t, err)
	})

	t.Run("should return error when sensor id is invalid", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		err := r.InsertSensorData(context.Background(), input.LatestData, "")

		// then
		assert.Error(t, err)
	})
}

var inputPayload = &entities.MqttPayload{
	Device:      "sensor-123",
	Battery:     34.0,
	Humidity:    50,
	Temperature: 20,
	Watermarks: []entities.Watermark{
		{
			Resistance: 23,
			Centibar:   38,
			Depth:      30,
		},
		{
			Resistance: 23,
			Centibar:   38,
			Depth:      60,
		},
		{
			Resistance: 23,
			Centibar:   38,
			Depth:      90,
		},
	},
}

var input = &entities.SensorCreate{
	ID:     "sensor-123",
	Status: entities.SensorStatusOnline,
	LatestData: &entities.SensorData{
		ID:        1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Data:      inputPayload,
	},
	Latitude:  9.446741,
	Longitude: 54.801539,
}

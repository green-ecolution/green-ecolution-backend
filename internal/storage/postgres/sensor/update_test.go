package sensor

import (
	"context"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestSensorRepository_Update(t *testing.T) {
	suite.ResetDB(t)
	suite.InsertSeed(t, "internal/storage/postgres/seed/test/sensor")

	t.Run("should update sensor successfully", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())
		newLat := 54.82078826498143
		newLong := 9.489684366114483

		got, err := r.Update(context.Background(),
			"sensor-1",
			WithStatus(entities.SensorStatusOffline),
			WithLatitude(newLat),
			WithLongitude(newLong))

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, entities.SensorStatusOffline, got.Status)
		assert.Equal(t, newLat, got.Latitude)
		assert.Equal(t, newLong, got.Longitude)
	})

	t.Run("should return error when update sensor with empty name", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.Update(context.Background(), "sensor-1", WithStatus(""))

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when update sensor with empty id", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.Update(context.Background(), "")

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when sensor not found", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.Update(context.Background(), "notFoundID")

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.Update(ctx, "sensor-1", WithStatus(entities.SensorStatusOffline))

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

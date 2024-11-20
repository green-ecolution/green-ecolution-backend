package sensor

import (
	"context"
	"testing"

	sensorUtils "github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/sensor"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestSensorRepository_Update(t *testing.T) {
	suite.ResetDB(t)
	suite.InsertSeed(t, "internal/storage/postgres/seed/test/sensor")

	t.Run("should update sensor successfully", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		got, err := r.Update(context.Background(), sensorUtils.TestSensorID, WithStatus(entities.SensorStatusOffline))
		gotByID, _ := r.GetByID(context.Background(), sensorUtils.TestSensorID)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, entities.SensorStatusOffline, got.Status)
		assert.Equal(t, entities.SensorStatusOffline, gotByID.Status)
	})

	t.Run("should return error when update sensor with empty name", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.Update(context.Background(), sensorUtils.TestSensorID, WithStatus(""))

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
		got, err := r.Update(ctx, sensorUtils.TestSensorID, WithStatus(entities.SensorStatusOffline))

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

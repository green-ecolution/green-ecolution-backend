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

		got, err := r.Update(context.Background(), 1, WithStatus(entities.SensorStatusOffline))
		gotByID, _ := r.GetByID(context.Background(), 1)

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
		got, err := r.Update(context.Background(), 1, WithStatus(""))

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when update sensor with negative id", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.Update(context.Background(), -1, WithStatus(entities.SensorStatusOffline))

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when update sensor with zero id", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.Update(context.Background(), 0)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when sensor not found", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())
		invalidID := int32(999)

		// when
		got, err := r.Update(context.Background(), invalidID)

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
		got, err := r.Update(ctx, 1, WithStatus(entities.SensorStatusOffline))

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

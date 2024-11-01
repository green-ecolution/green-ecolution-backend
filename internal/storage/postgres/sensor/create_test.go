package sensor

import (
	"context"
	"testing"

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
}
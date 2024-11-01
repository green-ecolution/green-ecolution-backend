package sensor

import (
	"context"
	"testing"

	sensorUtils "github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/sensor"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	storageMock "github.com/green-ecolution/green-ecolution-backend/internal/storage/_mock"
	"github.com/stretchr/testify/assert"
)


func TestSensorService_GetAll(t *testing.T) {
	t.Run("should return all sensor", func(t *testing.T) {
		// given
		repo := storageMock.NewMockSensorRepository(t)
		svc := NewSensorService(repo)

		// when
		repo.EXPECT().GetAll(context.Background()).Return(sensorUtils.TestSensorList, nil)
		sensors, err := svc.GetAll(context.Background())

		// then
		assert.NoError(t, err)
		assert.Equal(t, sensorUtils.TestSensorList, sensors)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// given
		repo := storageMock.NewMockSensorRepository(t)
		svc := NewSensorService(repo)

		repo.EXPECT().GetAll(context.Background()).Return(nil, storage.ErrSensorNotFound)
		sensors, err := svc.GetAll(context.Background())

		// then
		assert.Error(t, err)
		assert.Nil(t, sensors)
	})
}

func TestReady(t *testing.T) {
	t.Run("should return true if the service is ready", func(t *testing.T) {
		// given
		repo := storageMock.NewMockSensorRepository(t)
		svc := NewSensorService(repo)

		// when
		ready := svc.Ready()

		// then
		assert.True(t, ready)
	})

	t.Run("should return false if the service is not ready", func(t *testing.T) {
		// given
		svc := NewSensorService(nil)

		// when
		ready := svc.Ready()

		// then
		assert.False(t, ready)
	})
}

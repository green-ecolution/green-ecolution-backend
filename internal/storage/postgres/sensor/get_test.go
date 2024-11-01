package sensor

import (
	"context"
	"fmt"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sensorUtils "github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/sensor"
	"github.com/stretchr/testify/assert"
)

func TestSensorRepository_GetAll(t *testing.T) {
	t.Run("should return all sensors", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/sensor")
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetAll(context.Background())

		fmt.Print(got[len(got)-1].Data)

		// then
		assert.NoError(t, err)
		assert.Equal(t, len(sensorUtils.TestSensorList), len(got))
		for i, sensor := range got {
			assert.Equal(t, sensorUtils.TestSensorList[i].ID, sensor.ID)
			assert.Equal(t, sensorUtils.TestSensorList[i].Status, sensor.Status)
			//assert.Equal(t, sensorUtils.TestSensorList[i].Data, sensor.Data)
			assert.NotZero(t, sensor.CreatedAt)
			assert.NotZero(t, sensor.UpdatedAt)
		}
	})

	t.Run("should return empty slice when db is empty", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetAll(context.Background())

		// then
		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.GetAll(ctx)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestSensorRepository_GetByID(t *testing.T) {
	t.Run("should return sensor by id", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/sensor")
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetByID(ctx, 1)

		// then
		assert.NoError(t, err)
		assert.Equal(t, sensorUtils.TestSensor.ID, got.ID)
		assert.Equal(t, sensorUtils.TestSensor.Status, got.Status)
		//assert.Equal(t, sensorUtils.TestSensor.Data, got.Data)
		assert.NotZero(t, got.CreatedAt)
		assert.NotZero(t, got.UpdatedAt)
	})

	t.Run("should return error when sensor not found", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetByID(ctx, 1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when sensor id is negative", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetByID(ctx, -1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when sensor id is zero", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetByID(ctx, 0)

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
		got, err := r.GetByID(ctx, 1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestSensorRepository_GetStatusByID(t *testing.T) {
	t.Run("should return sensor status by id", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/sensor")
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetStatusByID(ctx, 1)

		// then
		assert.NoError(t, err)
		assert.Equal(t, sensorUtils.TestSensor.Status, *got)
	})

	t.Run("should return error when sensor not found", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetStatusByID(ctx, 1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when sensor id is negative", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetStatusByID(ctx, -1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when sensor id is zero", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetStatusByID(ctx, 0)

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
		got, err := r.GetStatusByID(ctx, 1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestSensorRepository_GetSensorByStatus(t *testing.T) {
	t.Run("should return sensors by status", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/sensor")
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetSensorByStatus(ctx, &sensorUtils.TestSensor.Status)

		// then
		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, sensorUtils.TestSensor.ID, got[0].ID)
		assert.Equal(t, sensorUtils.TestSensor.Status, got[0].Status)
	})

	t.Run("should return empty slice when no sensors match status", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		status := entities.SensorStatus("offline")
		got, err := r.GetSensorByStatus(ctx, &status)

		// then
		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("should return error when status is nil", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetSensorByStatus(ctx, nil)

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
		got, err := r.GetSensorByStatus(ctx, &sensorUtils.TestSensor.Status)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestSensorRepository_GetSensorDataByID(t *testing.T) {
	t.Run("should return sensor data for valid id", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/sensor")
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetSensorDataByID(ctx, 1)

		// then
		assert.NoError(t, err)
		assert.NotEmpty(t, got)
		for _, data := range got {
			assert.Equal(t, int32(1), data.ID)
			assert.NotZero(t, data.CreatedAt)
			assert.NotZero(t, data.UpdatedAt)
		}
	})

	t.Run("should return empty slice when no data found", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetSensorDataByID(ctx, 999)

		// then
		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.GetSensorDataByID(ctx, 1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

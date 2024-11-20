package sensor

import (
	"context"
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

		// then
		assert.NoError(t, err)
		assert.Equal(t, len(sensorUtils.TestSensorList), len(got))
		for i, sensor := range got {
			assert.Equal(t, sensorUtils.TestSensorList[i].ID, sensor.ID)
			assert.Equal(t, sensorUtils.TestSensorList[i].Status, sensor.Status)
			assert.Equal(t, sensorUtils.TestSensorList[i].Latitude, sensor.Latitude)
			assert.Equal(t, sensorUtils.TestSensorList[i].Longitude, sensor.Longitude)
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
		got, err := r.GetByID(ctx, sensorUtils.TestSensorID)

		// then
		assert.NoError(t, err)
		assert.Equal(t, sensorUtils.TestSensor.ID, got.ID)
		assert.Equal(t, sensorUtils.TestSensor.Status, got.Status)
		assert.Equal(t, sensorUtils.TestSensor.Latitude, got.Latitude)
		assert.Equal(t, sensorUtils.TestSensor.Longitude, got.Longitude)
		assert.NotZero(t, got.CreatedAt)
		assert.NotZero(t, got.UpdatedAt)

	})

	t.Run("should return error when sensor not found", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetByID(ctx, sensorUtils.TestSensorID)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when sensor id is empty", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetByID(ctx, "")

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
		got, err := r.GetByID(ctx, "sensor-1")

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
		got, err := r.GetStatusByID(ctx, "sensor-1")

		// then
		assert.NoError(t, err)
		assert.Equal(t, entities.SensorStatusOnline, *got)
	})

	t.Run("should return error when sensor not found", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetStatusByID(ctx, sensorUtils.TestSensorID)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when sensor id is empty", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetStatusByID(ctx, "")

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())
		suite.ResetDB(t)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.GetStatusByID(ctx, sensorUtils.TestSensorID)

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

		testSensorList := []*entities.Sensor{sensorUtils.TestSensorList[0], sensorUtils.TestSensorList[3]}

		// then

		assert.NoError(t, err)
		assert.NotEmpty(t, got)
		assert.Len(t, got, len(testSensorList))
		for i := range got {
			assert.Equal(t, testSensorList[i].ID, got[i].ID)
			assert.Equal(t, testSensorList[i].Status, got[i].Status)
			assert.Equal(t, testSensorList[i].Latitude, got[i].Latitude)
			assert.Equal(t, testSensorList[i].Longitude, got[i].Longitude)
			assert.NotZero(t, got[i].CreatedAt)
			assert.NotZero(t, got[i].UpdatedAt)
		}
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
		got, err := r.GetSensorDataByID(ctx, sensorUtils.TestSensorID)

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
		got, err := r.GetSensorDataByID(ctx, "notFoundID")

		// then
		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())
		suite.ResetDB(t)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.GetSensorDataByID(ctx, sensorUtils.TestSensorID)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

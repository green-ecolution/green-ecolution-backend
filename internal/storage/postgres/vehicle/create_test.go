package vehicle

import (
	"context"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestVehicleRepository_Create(t *testing.T) {
	suite.ResetDB(t)
	suite.InsertSeed(t, "internal/storage/postgres/seed/test/vehicle")
	input := entities.Vehicle{
		Description:   "Big car",
		WaterCapacity: 2000,
		Type:          entities.VehicleTypeTrailer,
		Status:        entities.VehicleStatusNotAvailable,
	}

	t.Run("should create vehicle", func(t *testing.T) {
		// given
		r := NewVehicleRepository(defaultFields.store, defaultFields.VehicleMappers)

		numberPlate := "FL ZU 9876"

		// when
		got, err := r.Create(
			context.Background(),
			WithDescription(input.Description),
			WithNumberPlate(numberPlate),
			WithWaterCapacity(input.WaterCapacity),
			WithVehicleStatus(input.Status),
			WithVehicleType(input.Type),
		)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, input.Description, got.Description)
		assert.Equal(t, numberPlate, got.NumberPlate)
		assert.Equal(t, input.WaterCapacity, got.WaterCapacity)
		assert.Equal(t, input.Type, got.Type)
		assert.Equal(t, input.Status, got.Status)
	})

	t.Run("should create vehicle with no description, type and status", func(t *testing.T) {
		// given
		r := NewVehicleRepository(defaultFields.store, defaultFields.VehicleMappers)

		numberPlate := "FL ZB 9876"

		// when
		got, err := r.Create(
			context.Background(),
			WithNumberPlate(numberPlate),
			WithWaterCapacity(input.WaterCapacity),
		)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, "", got.Description)
		assert.Equal(t, numberPlate, got.NumberPlate)
		assert.Equal(t, input.WaterCapacity, got.WaterCapacity)
		assert.Equal(t, entities.VehicleTypeUnknown, got.Type)
		assert.Equal(t, entities.VehicleStatusUnknown, got.Status)
	})

	t.Run("should return error when create vehicle with zero water capacity", func(t *testing.T) {
		// given
		r := NewVehicleRepository(defaultFields.store, defaultFields.VehicleMappers)

		// when
		got, err := r.Create(
			context.Background(),
			WithDescription(input.Description),
			WithNumberPlate(input.NumberPlate),
			WithWaterCapacity(0),
		)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
	t.Run("should return error when create vehicle with no number plate", func(t *testing.T) {
		// given
		r := NewVehicleRepository(defaultFields.store, defaultFields.VehicleMappers)

		// when
		got, err := r.Create(
			context.Background(),
			WithDescription(input.Description),
			WithNumberPlate(""),
			WithWaterCapacity(input.WaterCapacity),
		)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when create vehicle with duplicate plate", func(t *testing.T) {
		// given
		r := NewVehicleRepository(defaultFields.store, defaultFields.VehicleMappers)

		numberPlate := "FL ZT 9876"

		// when
		firstVehicle, err := r.Create(
			context.Background(),
			WithDescription(input.Description),
			WithNumberPlate(numberPlate),
			WithWaterCapacity(input.WaterCapacity),
			WithVehicleStatus(input.Status),
			WithVehicleType(input.Type),
		)
		assert.NoError(t, err)
		assert.NotNil(t, firstVehicle)

		secondVehicle, err := r.Create(
			context.Background(),
			WithDescription("New Car"),
			WithNumberPlate(numberPlate),
			WithWaterCapacity(2.000),
		)

		assert.Error(t, err)
		assert.Nil(t, secondVehicle)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewVehicleRepository(defaultFields.store, defaultFields.VehicleMappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.Create(
			ctx,
			WithDescription(input.Description),
			WithNumberPlate(input.NumberPlate),
			WithWaterCapacity(input.WaterCapacity),
		)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

package vehicle

import (
	"context"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestVehicleRepository_UpdateSuite(t *testing.T) {
	suite.ResetDB(t)
	suite.InsertSeed(t, "internal/storage/postgres/seed/test/vehicle")

	input := entities.Vehicle{
		Description:   "Updated description",
		NumberPlate:   "FL NEW 9876",
		WaterCapacity: 10000,
		Type:          entities.VehicleTypeTransporter,
		Status:        entities.VehicleStatusAvailable,
	}

	t.Run("should update vehicle", func(t *testing.T) {
		// given
		r := NewVehicleRepository(defaultFields.store, defaultFields.VehicleMappers)

		// when
		got, err := r.Update(context.Background(),
			1,
			WithDescription(input.Description),
			WithNumberPlate(input.NumberPlate),
			WithWaterCapacity(input.WaterCapacity),
			WithVehicleStatus(input.Status),
			WithVehicleType(input.Type),
		)
		gotByID, _ := r.GetByID(context.Background(), 1)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, input.Description, got.Description)
		assert.Equal(t, input.NumberPlate, got.NumberPlate)
		assert.Equal(t, input.WaterCapacity, got.WaterCapacity)
		assert.Equal(t, input.Description, gotByID.Description)
		assert.Equal(t, input.NumberPlate, gotByID.NumberPlate)
		assert.Equal(t, input.WaterCapacity, gotByID.WaterCapacity)
		assert.Equal(t, input.Type, gotByID.Type)
		assert.Equal(t, input.Status, gotByID.Status)
	})

	t.Run("should return error when update vehicle with duplicate plate", func(t *testing.T) {
		// given
		r := NewVehicleRepository(defaultFields.store, defaultFields.VehicleMappers)

		numberPlate :=  "FL ZT 9876"

		// when
		firstVehicle, err := r.Update(
			context.Background(),
			int32(2),
			WithDescription(input.Description),
			WithNumberPlate(numberPlate),
			WithWaterCapacity(input.WaterCapacity),
			WithVehicleStatus(input.Status),
			WithVehicleType(input.Type),
		)
		assert.NoError(t, err)
		assert.NotNil(t, firstVehicle)

		secondVehicle, err := r.Update(
			context.Background(),
			int32(1),
			WithDescription("New Car"),
			WithNumberPlate(numberPlate),
			WithWaterCapacity(2.000),
		)

		assert.Error(t, err)
		assert.Nil(t, secondVehicle)
	})

	t.Run("should return error when update vehicle with zero water capacity", func(t *testing.T) {
		// given
		r := NewVehicleRepository(defaultFields.store, defaultFields.VehicleMappers)
		// when
		got, err := r.Update(context.Background(), 2, WithWaterCapacity(0))

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when update vehicle with no number plate", func(t *testing.T) {
		// given
		r := NewVehicleRepository(defaultFields.store, defaultFields.VehicleMappers)
		// when
		got, err := r.Update(context.Background(), 2, WithNumberPlate(""))

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when update vehicle with negative id", func(t *testing.T) {
		// given
		r := NewVehicleRepository(defaultFields.store, defaultFields.VehicleMappers)

		// when
		got, err := r.Update(context.Background(), -1, WithNumberPlate(input.NumberPlate))

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when update vehicle with zero id", func(t *testing.T) {
		// given
		r := NewVehicleRepository(defaultFields.store, defaultFields.VehicleMappers)

		// when
		got, err := r.Update(context.Background(), 0, WithNumberPlate(input.NumberPlate))

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when update vehicle not found", func(t *testing.T) {
		// given
		r := NewVehicleRepository(defaultFields.store, defaultFields.VehicleMappers)

		// when
		got, err := r.Update(context.Background(), 99, WithNumberPlate(input.NumberPlate))

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewVehicleRepository(defaultFields.store, defaultFields.VehicleMappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.Update(ctx, 1, WithNumberPlate(input.NumberPlate))

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

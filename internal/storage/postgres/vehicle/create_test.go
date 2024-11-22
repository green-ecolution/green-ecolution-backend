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
		DriverLicense: entities.DriverLicenseTrailer,
		Height:        1.5,
		Length:        2.0,
		Width:         2.0,
		Model:         "1615/17 - Conrad - MAN TGE 3.180",
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
			WithModel(input.Model),
			WithDriverLicense(input.DriverLicense),
			WithHeight(input.Height),
			WithLength(input.Length),
			WithWidth(input.Width),
		)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, input.Description, got.Description)
		assert.Equal(t, numberPlate, got.NumberPlate)
		assert.Equal(t, input.WaterCapacity, got.WaterCapacity)
		assert.Equal(t, input.Type, got.Type)
		assert.Equal(t, input.Status, got.Status)
		assert.Equal(t, input.DriverLicense, got.DriverLicense)
		assert.Equal(t, input.Height, got.Height)
		assert.Equal(t, input.Length, got.Length)
		assert.Equal(t, input.Width, got.Width)
		assert.Equal(t, input.Model, got.Model)
	})

	t.Run("should create vehicle with no description, type, model, driver license and status", func(t *testing.T) {
		// given
		r := NewVehicleRepository(defaultFields.store, defaultFields.VehicleMappers)

		numberPlate := "FL ZB 9876"

		// when
		got, err := r.Create(
			context.Background(),
			WithNumberPlate(numberPlate),
			WithWaterCapacity(input.WaterCapacity),
			WithHeight(input.Height),
			WithLength(input.Length),
			WithWidth(input.Width),
		)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, "", got.Description)
		assert.Equal(t, "", got.Model)
		assert.Equal(t, numberPlate, got.NumberPlate)
		assert.Equal(t, input.WaterCapacity, got.WaterCapacity)
		assert.Equal(t, entities.VehicleTypeUnknown, got.Type)
		assert.Equal(t, entities.VehicleStatusUnknown, got.Status)
		assert.Equal(t, entities.DriverLicenseCar, got.DriverLicense)
		assert.Equal(t, input.Height, got.Height)
		assert.Equal(t, input.Length, got.Length)
		assert.Equal(t, input.Width, got.Width)
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
			WithHeight(input.Height),
			WithLength(input.Length),
			WithWidth(input.Width),
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
			WithHeight(input.Height),
			WithLength(input.Length),
			WithWidth(input.Width),
		)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when create vehicle with zero size measurements", func(t *testing.T) {
		// given
		r := NewVehicleRepository(defaultFields.store, defaultFields.VehicleMappers)

		numberPlate := "FL ZB 9876"

		// when
		got, err := r.Create(
			context.Background(),
			WithDescription(input.Description),
			WithNumberPlate(numberPlate),
			WithWaterCapacity(input.WaterCapacity),
			WithHeight(0),
			WithLength(0),
			WithWidth(0),
		)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when create vehicle with wrong driver license", func(t *testing.T) {
		// given
		r := NewVehicleRepository(defaultFields.store, defaultFields.VehicleMappers)

		// when
		got, err := r.Create(
			context.Background(),
			WithDescription(input.Description),
			WithNumberPlate(input.NumberPlate),
			WithWaterCapacity(input.WaterCapacity),
			WithDriverLicense(""),
			WithHeight(input.Height),
			WithLength(input.Length),
			WithWidth(input.Width),
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
			WithHeight(input.Height),
			WithLength(input.Length),
			WithWidth(input.Width),
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

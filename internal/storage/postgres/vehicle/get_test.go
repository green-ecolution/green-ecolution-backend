package vehicle

import (
	"context"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestVehicleRepository_GetAll(t *testing.T) {
	t.Run("should return all verhicles", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/vehicle")
		r := NewVehicleRepository(suite.Store, defaultVehicleMappers())

		// when
		got, err := r.GetAll(context.Background())

		// then
		assert.NoError(t, err)
		assert.Equal(t, len(allTestVehicles), len(got))
		for i, vehicle := range got {
			assert.Equal(t, allTestVehicles[i].ID, vehicle.ID)
			assert.Equal(t, allTestVehicles[i].Description, vehicle.Description)
			assert.Equal(t, allTestVehicles[i].NumberPlate, vehicle.NumberPlate)
			assert.Equal(t, allTestVehicles[i].WaterCapacity, vehicle.WaterCapacity)
			assert.Equal(t, allTestVehicles[i].Type, vehicle.Type)
			assert.Equal(t, allTestVehicles[i].Status, vehicle.Status)
			assert.Equal(t, allTestVehicles[i].DrivingLicense, vehicle.DrivingLicense)
			assert.Equal(t, allTestVehicles[i].Height, vehicle.Height)
			assert.Equal(t, allTestVehicles[i].Width, vehicle.Width)
			assert.Equal(t, allTestVehicles[i].Length, vehicle.Length)
			assert.Equal(t, allTestVehicles[i].Model, vehicle.Model)
			assert.NotZero(t, vehicle.CreatedAt)
			assert.NotZero(t, vehicle.UpdatedAt)
		}
	})

	t.Run("should return empty slice when db is empty", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewVehicleRepository(suite.Store, defaultVehicleMappers())

		// when
		got, err := r.GetAll(context.Background())

		// then
		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewVehicleRepository(suite.Store, defaultVehicleMappers())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.GetAll(ctx)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestVehicleRepository_GetAllByType(t *testing.T) {
	t.Run("should return all verhicles of type transporter", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/vehicle")
		r := NewVehicleRepository(suite.Store, defaultVehicleMappers())

		// when
		got, err := r.GetAllByType(context.Background(), entities.VehicleTypeTransporter)

		// then
		assert.NoError(t, err)
		assert.Equal(t, 1, len(got))
		for _, vehicle := range got {
			assert.Equal(t, entities.VehicleTypeTransporter, vehicle.Type)
		}
	})

	t.Run("should return all verhicles of type trailer", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/vehicle")
		r := NewVehicleRepository(suite.Store, defaultVehicleMappers())

		// when
		got, err := r.GetAllByType(context.Background(), entities.VehicleTypeTrailer)

		// then
		assert.NoError(t, err)
		assert.Equal(t, 1, len(got))
		for _, vehicle := range got {
			assert.Equal(t, entities.VehicleTypeTrailer, vehicle.Type)
		}
	})

	t.Run("should return empty slice when db is empty", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewVehicleRepository(suite.Store, defaultVehicleMappers())

		// when
		got, err := r.GetAllByType(context.Background(), entities.VehicleTypeUnknown)

		// then
		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewVehicleRepository(suite.Store, defaultVehicleMappers())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.GetAllByType(ctx, entities.VehicleTypeUnknown)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestVehicleRepository_GetByID(t *testing.T) {
	t.Run("should return verhicle by id", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/vehicle")
		r := NewVehicleRepository(suite.Store, defaultVehicleMappers())
		shouldReturn := allTestVehicles[0]

		// when
		got, err := r.GetByID(ctx, 1)

		// then
		assert.NoError(t, err)
		assert.Equal(t, shouldReturn.ID, got.ID)
		assert.Equal(t, shouldReturn.NumberPlate, got.NumberPlate)
		assert.Equal(t, shouldReturn.Description, got.Description)
		assert.Equal(t, shouldReturn.WaterCapacity, got.WaterCapacity)
		assert.Equal(t, shouldReturn.Type, got.Type)
		assert.Equal(t, shouldReturn.Status, got.Status)
		assert.Equal(t, shouldReturn.DrivingLicense, got.DrivingLicense)
		assert.Equal(t, shouldReturn.Height, got.Height)
		assert.Equal(t, shouldReturn.Length, got.Length)
		assert.Equal(t, shouldReturn.Width, got.Width)
		assert.Equal(t, shouldReturn.Model, got.Model)
		assert.NotZero(t, got.CreatedAt)
		assert.NotZero(t, got.UpdatedAt)
	})

	t.Run("should return error when verhicle not found", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewVehicleRepository(suite.Store, defaultVehicleMappers())

		// when
		got, err := r.GetByID(ctx, 1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when vehicle id is negative", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewVehicleRepository(suite.Store, defaultVehicleMappers())

		// when
		got, err := r.GetByID(ctx, -1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when vehicle id is zero", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewVehicleRepository(suite.Store, defaultVehicleMappers())

		// when
		got, err := r.GetByID(ctx, 0)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewVehicleRepository(suite.Store, defaultVehicleMappers())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.GetByID(ctx, 1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestVehicleRepository_GetByPlate(t *testing.T) {
	tests := []struct {
		name string
		want *entities.Vehicle
		args string
	}{
		{
			name: "should return region by plate 'B-1234'",
			want: allTestVehicles[0],
			args: "B-1234",
		},
		{
			name: "should return region by name 'B-5678'",
			want: allTestVehicles[1],
			args: "B-5678",
		},
	}

	suite.ResetDB(t)
	suite.InsertSeed(t, "internal/storage/postgres/seed/test/vehicle")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ctx := context.Background()
			r := NewVehicleRepository(suite.Store, defaultVehicleMappers())

			// when
			got, err := r.GetByPlate(ctx, tt.args)

			// then
			assert.NoError(t, err)
			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.NumberPlate, got.NumberPlate)
			assert.Equal(t, tt.want.Description, got.Description)
			assert.Equal(t, tt.want.WaterCapacity, got.WaterCapacity)
			assert.Equal(t, tt.want.Type, got.Type)
			assert.Equal(t, tt.want.Status, got.Status)
			assert.Equal(t, tt.want.DrivingLicense, got.DrivingLicense)
			assert.Equal(t, tt.want.Height, got.Height)
			assert.Equal(t, tt.want.Length, got.Length)
			assert.Equal(t, tt.want.Width, got.Width)
			assert.Equal(t, tt.want.Model, got.Model)
			assert.NotZero(t, got.CreatedAt)
			assert.NotZero(t, got.UpdatedAt)
		})
	}

	t.Run("should return error when vehicle not found", func(t *testing.T) {
		// given
		ctx := context.Background()
		r := NewVehicleRepository(suite.Store, defaultVehicleMappers())

		// when
		got, err := r.GetByPlate(ctx, "Non-existing")

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when vehicle plate is empty", func(t *testing.T) {
		// given
		ctx := context.Background()
		r := NewVehicleRepository(suite.Store, defaultVehicleMappers())

		// when
		got, err := r.GetByPlate(ctx, "")

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewVehicleRepository(suite.Store, defaultVehicleMappers())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.GetByPlate(ctx, "B-1234")

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

var allTestVehicles = []*entities.Vehicle{
	{
		ID:             1,
		NumberPlate:    "B-1234",
		Description:    "Test vehicle 1",
		WaterCapacity:  100.0,
		Type:           entities.VehicleTypeTrailer,
		Status:         entities.VehicleStatusActive,
		Model:          "1615/17 - Conrad - MAN TGE 3.180",
		DrivingLicense: entities.DrivingLicenseTrailer,
		Height:         1.5,
		Length:         2.0,
		Width:          2.0,
	},
	{
		ID:             2,
		NumberPlate:    "B-5678",
		Description:    "Test vehicle 2",
		WaterCapacity:  150.0,
		Type:           entities.VehicleTypeTransporter,
		Status:         entities.VehicleStatusUnknown,
		Model:          "Actros L Mercedes Benz",
		DrivingLicense: entities.DrivingLicenseTransporter,
		Height:         2.1,
		Length:         5.0,
		Width:          2.4,
	},
}

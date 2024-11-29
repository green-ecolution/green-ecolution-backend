package wateringplan

import (
	"context"
	"testing"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestWateringPlanRepository_GetAll(t *testing.T) {
	t.Run("should return all watering plans", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/watering_plan")
		r := NewWateringPlanRepository(suite.Store, mappers)

		// when
		got, err := r.GetAll(context.Background())

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.NotEmpty(t, got)
		assert.Len(t, got, len(allTestWateringPlans))
		for i, wp := range got {
			assert.Equal(t, allTestWateringPlans[i].ID, wp.ID)
			assert.Equal(t, allTestWateringPlans[i].Date, wp.Date)
			assert.Equal(t, allTestWateringPlans[i].Description, wp.Description)
			assert.Equal(t, allTestWateringPlans[i].WateringPlanStatus, wp.WateringPlanStatus)
			assert.Equal(t, allTestWateringPlans[i].Distance, wp.Distance)
			assert.Equal(t, allTestWateringPlans[i].TotalWaterRequired, wp.TotalWaterRequired)

			// assert transporter
			assert.Equal(t, allTestWateringPlans[i].Transporter.ID, wp.Transporter.ID)

			// assert trailer
			if allTestWateringPlans[i].Trailer == nil {
				assert.Nil(t, wp.Trailer)
				assert.NoError(t, err)
			} else {
				assert.NotNil(t, wp.Trailer)
				assert.Equal(t, allTestWateringPlans[i].Trailer.ID, wp.Trailer.ID)
			}
		}
	})

	t.Run("should return empty slice when db is empty", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewWateringPlanRepository(suite.Store, mappers)

		// when
		got, err := r.GetAll(context.Background())

		// then
		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		_, err := r.GetAll(ctx)

		// then
		assert.Error(t, err)
	})
}

func TestTreeClusterRepository_GetByID(t *testing.T) {
	suite.ResetDB(t)
	suite.InsertSeed(t, "internal/storage/postgres/seed/test/watering_plan")

	t.Run("should return watering plan by id", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		// when
		got, err := r.GetByID(context.Background(), 1)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, allTestWateringPlans[0].ID, got.ID)
		assert.Equal(t, allTestWateringPlans[0].Date, got.Date)
		assert.Equal(t, allTestWateringPlans[0].Description, got.Description)
		assert.Equal(t, allTestWateringPlans[0].WateringPlanStatus, got.WateringPlanStatus)
		assert.Equal(t, allTestWateringPlans[0].Distance, got.Distance)
		assert.Equal(t, allTestWateringPlans[0].TotalWaterRequired, got.TotalWaterRequired)

		// assert transporter
		assert.Equal(t, allTestWateringPlans[0].Transporter.ID, got.Transporter.ID)
		assert.Equal(t, allTestWateringPlans[0].Transporter.Type, got.Transporter.Type)

		// assert trailer
		assert.Equal(t, allTestWateringPlans[0].Trailer.ID, got.Trailer.ID)
		assert.Equal(t, allTestWateringPlans[0].Trailer.Type, got.Trailer.Type)
	})

	t.Run("should return watering plan by id without trailer", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		// when
		got, err := r.GetByID(context.Background(), 2)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, allTestWateringPlans[1].ID, got.ID)
		assert.Equal(t, allTestWateringPlans[1].Date, got.Date)
		assert.Equal(t, allTestWateringPlans[1].Description, got.Description)
		assert.Equal(t, allTestWateringPlans[1].WateringPlanStatus, got.WateringPlanStatus)
		assert.Equal(t, allTestWateringPlans[1].Distance, got.Distance)
		assert.Equal(t, allTestWateringPlans[1].TotalWaterRequired, got.TotalWaterRequired)

		// assert transporter
		assert.Equal(t, allTestWateringPlans[1].Transporter.ID, got.Transporter.ID)
		assert.Equal(t, allTestWateringPlans[1].Transporter.Type, got.Transporter.Type)

		// assert nil trailer
		assert.Nil(t, got.Trailer)
	})

	t.Run("should return error when watering plan with non-existing id", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		// when
		got, err := r.GetByID(context.Background(), 99)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when watering plan with negative id", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		// when
		got, err := r.GetByID(context.Background(), -1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when watering plan with zero id", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		// when
		got, err := r.GetByID(context.Background(), 0)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.GetByID(ctx, 1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestWateringPlanRepository_GetLinkedVehicleByID(t *testing.T) {
	ctx := context.Background()
	suite.ResetDB(t)
	suite.InsertSeed(t, "internal/storage/postgres/seed/test/watering_plan")

	t.Run("should return vehicle with type transporter by watering plan id", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)
		shouldReturn := allTestVehicles[1]

		// when
		got, err := r.GetLinkedVehicleByID(ctx, int32(1), entities.VehicleTypeTransporter)

		// then
		assert.NoError(t, err)
		assert.Equal(t, shouldReturn.ID, got.ID)
		assert.Equal(t, shouldReturn.NumberPlate, got.NumberPlate)
		assert.Equal(t, shouldReturn.Description, got.Description)
		assert.Equal(t, shouldReturn.WaterCapacity, got.WaterCapacity)
		assert.Equal(t, entities.VehicleTypeTransporter, got.Type)
		assert.Equal(t, shouldReturn.Status, got.Status)
		assert.NotZero(t, got.CreatedAt)
		assert.NotZero(t, got.UpdatedAt)
	})

	t.Run("should return vehicle with type trailer by watering plan id", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)
		shouldReturn := allTestVehicles[0]

		// when
		got, err := r.GetLinkedVehicleByID(ctx, int32(1), entities.VehicleTypeTrailer)

		// then
		assert.NoError(t, err)
		assert.Equal(t, shouldReturn.ID, got.ID)
		assert.Equal(t, shouldReturn.NumberPlate, got.NumberPlate)
		assert.Equal(t, shouldReturn.Description, got.Description)
		assert.Equal(t, shouldReturn.WaterCapacity, got.WaterCapacity)
		assert.Equal(t, entities.VehicleTypeTrailer, got.Type)
		assert.Equal(t, shouldReturn.Status, got.Status)
		assert.NotZero(t, got.CreatedAt)
		assert.NotZero(t, got.UpdatedAt)
	})

	t.Run("should return error when watering plan not found", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		// when
		got, err := r.GetLinkedVehicleByID(ctx, int32(99), entities.VehicleTypeTrailer)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when vehicle with type trailer is not found", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		// when
		got, err := r.GetLinkedVehicleByID(ctx, int32(2), entities.VehicleTypeTrailer)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when watering plan id is negative", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		// when
		got, err := r.GetLinkedVehicleByID(ctx, int32(-1), entities.VehicleTypeTransporter)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when watering plan id is zero", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		// when
		got, err := r.GetLinkedVehicleByID(ctx, int32(0), entities.VehicleTypeTransporter)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when vehicle type is not trailer or transporter", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		// when
		got, err := r.GetLinkedVehicleByID(ctx, int32(1), entities.VehicleTypeUnknown)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.GetLinkedVehicleByID(ctx, int32(1), entities.VehicleTypeTransporter)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

var allTestWateringPlans = []*entities.WateringPlan{
	{
		ID:                 1,
		Date:               time.Date(2024, 9, 22, 0, 0, 0, 0, time.UTC),
		Description:        "New watering plan for the west side of the city",
		WateringPlanStatus: entities.WateringPlanStatusPlanned,
		Distance:           utils.P(63.0),
		TotalWaterRequired: utils.P(6000.0),
		Transporter:        allTestVehicles[1],
		Trailer:            allTestVehicles[0],
	},
	{
		ID:                 2,
		Date:               time.Date(2024, 8, 3, 0, 0, 0, 0, time.UTC),
		Description:        "New watering plan for the east side of the city",
		WateringPlanStatus: entities.WateringPlanStatusActive,
		Distance:           utils.P(63.0),
		TotalWaterRequired: utils.P(6000.0),
		Transporter:        allTestVehicles[1],
		Trailer:            nil,
	},
	{
		ID:                 3,
		Date:               time.Date(2024, 6, 12, 0, 0, 0, 0, time.UTC),
		Description:        "Very important watering plan due to no rainfall",
		WateringPlanStatus: entities.WateringPlanStatusFinished,
		Distance:           utils.P(63.0),
		TotalWaterRequired: utils.P(6000.0),
		Transporter:        allTestVehicles[1],
		Trailer:            nil,
	},
	{
		ID:                 4,
		Date:               time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC),
		Description:        "New watering plan for the south side of the city",
		WateringPlanStatus: entities.WateringPlanStatusNotCompeted,
		Distance:           utils.P(63.0),
		TotalWaterRequired: utils.P(6000.0),
		Transporter:        allTestVehicles[1],
		Trailer:            nil,
	},
	{
		ID:                 5,
		Date:               time.Date(2024, 6, 4, 0, 0, 0, 0, time.UTC),
		Description:        "Canceled due to flood",
		WateringPlanStatus: entities.WateringPlanStatusCanceled,
		Distance:           utils.P(63.0),
		TotalWaterRequired: utils.P(6000.0),
		Transporter:        allTestVehicles[1],
		Trailer:            nil,
	},
}

var allTestVehicles = []*entities.Vehicle{
	{
		ID:            1,
		NumberPlate:   "B-1234",
		Description:   "Test vehicle 1",
		WaterCapacity: 100.0,
		Type:          entities.VehicleTypeTrailer,
		Status:        entities.VehicleStatusActive,
	},
	{
		ID:            2,
		NumberPlate:   "B-5678",
		Description:   "Test vehicle 2",
		WaterCapacity: 150.0,
		Type:          entities.VehicleTypeTransporter,
		Status:        entities.VehicleStatusUnknown,
	},
}

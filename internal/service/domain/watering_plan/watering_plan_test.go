package wateringplan

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	storageMock "github.com/green-ecolution/green-ecolution-backend/internal/storage/_mock"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestWateringPlanService_GetAll(t *testing.T) {
	ctx := context.Background()

	t.Run("should return all watering plans when successful", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo)

		wateringPlanRepo.EXPECT().GetAll(ctx).Return(allTestWateringPlans, nil)

		// when
		wateringPlans, err := svc.GetAll(ctx)

		// then
		assert.NoError(t, err)
		assert.Equal(t, allTestWateringPlans, wateringPlans)
	})

	t.Run("should return empty slice when no watering plans are found", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo)

		wateringPlanRepo.EXPECT().GetAll(ctx).Return([]*entities.WateringPlan{}, nil)

		// when
		wateringPlans, err := svc.GetAll(ctx)

		// then
		assert.NoError(t, err)
		assert.Empty(t, wateringPlans)
	})

	t.Run("should return error when GetAll fails", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo)

		expectedErr := errors.New("GetAll failed")
		wateringPlanRepo.EXPECT().GetAll(ctx).Return(nil, expectedErr)

		// when
		wateringPlans, err := svc.GetAll(ctx)

		// then
		assert.Error(t, err)
		assert.Nil(t, wateringPlans)
		assert.Equal(t, "500: GetAll failed", err.Error())
	})
}

func TestWateringPlanService_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("should return watering plan when found", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo)

		id := int32(1)
		expectedPlan := allTestWateringPlans[0]
		wateringPlanRepo.EXPECT().GetByID(ctx, id).Return(expectedPlan, nil)

		// when
		wateringPlan, err := svc.GetByID(ctx, id)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedPlan, wateringPlan)
	})

	t.Run("should return error if watering plan not found", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo)

		id := int32(1)
		expectedErr := storage.ErrEntityNotFound
		wateringPlanRepo.EXPECT().GetByID(ctx, id).Return(nil, expectedErr)

		// when
		wateringPlan, err := svc.GetByID(ctx, id)

		// then
		assert.Error(t, err)
		assert.Nil(t, wateringPlan)
		assert.Equal(t, "404: watering plan not found", err.Error())
	})
}

var allTestWateringPlans = []*entities.WateringPlan{
	{
		ID:                 1,
		Date:               time.Date(2024, 9, 22, 0, 0, 0, 0, time.UTC),
		Description:        "New watering plan for the west side of the city",
		WateringPlanStatus: "planned",
		Distance:           utils.P(63.0),
		TotalWaterRequired: utils.P(6000.0),
	},
	{
		ID:                 2,
		Date:               time.Date(2024, 8, 3, 0, 0, 0, 0, time.UTC),
		Description:        "New watering plan for the east side of the city",
		WateringPlanStatus: "active",
		Distance:           utils.P(63.0),
		TotalWaterRequired: utils.P(6000.0),
	},
	{
		ID:                 3,
		Date:               time.Date(2024, 6, 12, 0, 0, 0, 0, time.UTC),
		Description:        "Very important watering plan due to no rainfall",
		WateringPlanStatus: "finished",
		Distance:           utils.P(63.0),
		TotalWaterRequired: utils.P(6000.0),
	},
	{
		ID:                 4,
		Date:               time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC),
		Description:        "New watering plan for the south side of the city",
		WateringPlanStatus: "not competed",
		Distance:           utils.P(63.0),
		TotalWaterRequired: utils.P(6000.0),
	},
	{
		ID:                 5,
		Date:               time.Date(2024, 6, 4, 0, 0, 0, 0, time.UTC),
		Description:        "Canceled due to flood",
		WateringPlanStatus: "canceled",
		Distance:           utils.P(63.0),
		TotalWaterRequired: utils.P(6000.0),
	},
}

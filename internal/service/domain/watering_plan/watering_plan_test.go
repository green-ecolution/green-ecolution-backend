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
	"github.com/stretchr/testify/mock"
)

func TestWateringPlanService_GetAll(t *testing.T) {
	ctx := context.Background()

	t.Run("should return all watering plans when successful", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo)

		wateringPlanRepo.EXPECT().GetAll(ctx).Return(allTestWateringPlans, nil)

		// when
		wateringPlans, err := svc.GetAll(ctx)

		// then
		assert.NoError(t, err)
		assert.Equal(t, allTestWateringPlans, wateringPlans)
	})

	t.Run("should return empty slice when no watering plans are found", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo)

		wateringPlanRepo.EXPECT().GetAll(ctx).Return([]*entities.WateringPlan{}, nil)

		// when
		wateringPlans, err := svc.GetAll(ctx)

		// then
		assert.NoError(t, err)
		assert.Empty(t, wateringPlans)
	})

	t.Run("should return error when GetAll fails", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo)

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
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo)

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
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo)

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

func TestWateringPlanService_Create(t *testing.T) {
	ctx := context.Background()
	newWateringPlan := &entities.WateringPlanCreate{
		Date:          time.Date(2024, 9, 26, 0, 0, 0, 0, time.UTC),
		Description:   "New watering plan",
		TransporterID: utils.P(int32(2)),
		TrailerID:     utils.P(int32(1)),
		TreeclusterIDs:       []*int32{utils.P(int32(1)), utils.P(int32(2))},
	}

	t.Run("should successfully create a new watering plan", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo)

		// check treecluster
		clusterRepo.EXPECT().GetByIDs(
			ctx,
			[]int32{1, 2},
		).Return(allTestClusters[0:2], nil)

		// check transporter
		vehicleRepo.EXPECT().GetByID(
			ctx,
			int32(2),
		).Return(allTestVehicles[1], nil)

		// check trailer
		vehicleRepo.EXPECT().GetByID(
			ctx,
			int32(1),
		).Return(allTestVehicles[0], nil)

		wateringPlanRepo.EXPECT().Create(
			ctx,
			mock.Anything,
		).Return(allTestWateringPlans[0], nil)

		// when
		result, err := svc.Create(ctx, newWateringPlan)

		// then
		assert.NoError(t, err)
		assert.Equal(t, allTestWateringPlans[0], result)
	})

	t.Run("should successfully create a new watering plan without a trailer", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo)

		newWateringPlan := &entities.WateringPlanCreate{
			Date:          time.Date(2024, 9, 26, 0, 0, 0, 0, time.UTC),
			Description:   "New watering plan",
			TransporterID: utils.P(int32(2)),
			TreeclusterIDs:       []*int32{utils.P(int32(1)), utils.P(int32(2))},
		}

		// check treecluster
		clusterRepo.EXPECT().GetByIDs(
			ctx,
			[]int32{1, 2},
		).Return(allTestClusters[0:2], nil)

		// check transporter
		vehicleRepo.EXPECT().GetByID(
			ctx,
			int32(2),
		).Return(allTestVehicles[1], nil)

		wateringPlanRepo.EXPECT().Create(
			ctx,
			mock.Anything,
		).Return(allTestWateringPlans[0], nil)

		// when
		result, err := svc.Create(ctx, newWateringPlan)

		// then
		assert.NoError(t, err)
		assert.Equal(t, allTestWateringPlans[0], result)
	})

	t.Run("should return an error when finding treeclusters fails", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo)

		// check treecluster
		clusterRepo.EXPECT().GetByIDs(
			ctx,
			[]int32{1, 2},
		).Return(nil, storage.ErrConnectionClosed)

		// when
		result, err := svc.Create(ctx, newWateringPlan)

		// then
		assert.Nil(t, result)
		assert.EqualError(t, err, "500: connection is closed")
	})

	t.Run("should return an error when treecluster are empty", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo)

		// check treecluster
		clusterRepo.EXPECT().GetByIDs(
			ctx,
			[]int32{1, 2},
		).Return([]*entities.TreeCluster{}, nil)

		// when
		result, err := svc.Create(ctx, newWateringPlan)

		// then
		assert.Nil(t, result)
		assert.EqualError(t, err, "404: treecluster not found")
	})

	t.Run("should return an error when transporter is not found", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo)

		// check treecluster
		clusterRepo.EXPECT().GetByIDs(
			ctx,
			[]int32{1, 2},
		).Return(allTestClusters[0:2], nil)

		// check transporter
		vehicleRepo.EXPECT().GetByID(
			ctx,
			int32(2),
		).Return(nil, storage.ErrVehicleNotFound)

		// when
		result, err := svc.Create(ctx, newWateringPlan)

		// then
		assert.Nil(t, result)
		assert.EqualError(t, err, "404: vehicle not found")
	})

	t.Run("should return an error when creating watering plan fails", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo)

		expectedErr := errors.New("Failed to create watering plan")

		// check treecluster
		clusterRepo.EXPECT().GetByIDs(
			ctx,
			[]int32{1, 2},
		).Return(allTestClusters[0:2], nil)

		// check transporter
		vehicleRepo.EXPECT().GetByID(
			ctx,
			int32(2),
		).Return(allTestVehicles[1], nil)

		// check trailer
		vehicleRepo.EXPECT().GetByID(
			ctx,
			int32(1),
		).Return(allTestVehicles[0], nil)

		wateringPlanRepo.EXPECT().Create(
			ctx,
			mock.Anything,
		).Return(nil, expectedErr)

		// when
		result, err := svc.Create(ctx, newWateringPlan)

		// then
		assert.Nil(t, result)
		assert.EqualError(t, err, "500: Failed to create watering plan")
	})

	t.Run("should return validation error on empty date", func(t *testing.T) {
		// given
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo)

		newWateringPlan.Date = time.Time{}

		// when
		result, err := svc.Create(ctx, newWateringPlan)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validation error")
	})

	t.Run("should return validation error on empty transporter", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo)

		newWateringPlan := &entities.WateringPlanCreate{
			Date:        time.Date(2024, 9, 26, 0, 0, 0, 0, time.UTC),
			Description: "New watering plan",
		}

		// when
		result, err := svc.Create(ctx, newWateringPlan)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validation error")
	})

	t.Run("should return validation error on empty treeclusters", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo)

		newWateringPlan := &entities.WateringPlanCreate{
			Date:        time.Date(2024, 9, 26, 0, 0, 0, 0, time.UTC),
			Description: "New watering plan",
			TransporterID: utils.P(int32(2)),
		}

		// when
		result, err := svc.Create(ctx, newWateringPlan)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validation error")
	})
}

func TestWateringPlanService_Update(t *testing.T) {
	ctx := context.Background()
	wateringPlanID := int32(1)
	updatedWateringPlan := &entities.WateringPlanUpdate{
		Date:          time.Date(2024, 8, 3, 0, 0, 0, 0, time.UTC),
		Description:   "New watering plan for the east side of the city",
		TransporterID: utils.P(int32(2)),
		TrailerID:     utils.P(int32(1)),
		TreeclusterIDs: []*int32{utils.P(int32(1)), utils.P(int32(2))},
	}

	t.Run("should successfully update a watering plan", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo)

		// check treecluster
		clusterRepo.EXPECT().GetByIDs(
			ctx,
			[]int32{1, 2},
		).Return(allTestClusters[0:2], nil)

		// check transporter
		vehicleRepo.EXPECT().GetByID(
			ctx,
			int32(2),
		).Return(allTestVehicles[1], nil)

		// check trailer
		vehicleRepo.EXPECT().GetByID(
			ctx,
			int32(1),
		).Return(allTestVehicles[0], nil)

		wateringPlanRepo.EXPECT().Update(
			ctx,
			wateringPlanID,
			mock.Anything,
		).Return(nil)

		wateringPlanRepo.EXPECT().GetByID(
			ctx,
			wateringPlanID,
		).Return(allTestWateringPlans[1], nil)

		// when
		result, err := svc.Update(ctx, wateringPlanID, updatedWateringPlan)

		// then
		assert.NoError(t, err)
		assert.Equal(t, allTestWateringPlans[1], result)
	})

	t.Run("should return error update a watering plan without a trailer", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo)

		updatedWateringPlan := &entities.WateringPlanUpdate{
			Date:          time.Date(2024, 8, 3, 0, 0, 0, 0, time.UTC),
			Description:   "New watering plan for the east side of the city",
			TransporterID: utils.P(int32(2)),
			TreeclusterIDs: []*int32{utils.P(int32(1)), utils.P(int32(2))},
		}

		// check treecluster
		clusterRepo.EXPECT().GetByIDs(
			ctx,
			[]int32{1, 2},
		).Return(allTestClusters[0:2], nil)

		// check transporter
		vehicleRepo.EXPECT().GetByID(
			ctx,
			int32(2),
		).Return(allTestVehicles[1], nil)

		wateringPlanRepo.EXPECT().Update(
			ctx,
			wateringPlanID,
			mock.Anything,
		).Return(nil)

		wateringPlanRepo.EXPECT().GetByID(
			ctx,
			wateringPlanID,
		).Return(allTestWateringPlans[1], nil)

		// when
		result, err := svc.Update(ctx, wateringPlanID, updatedWateringPlan)

		// then
		assert.NoError(t, err)
		assert.Equal(t, allTestWateringPlans[1], result)
	})

	t.Run("should return an error when transporter is not found", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo)

		// check treecluster
		clusterRepo.EXPECT().GetByIDs(
			ctx,
			[]int32{1, 2},
		).Return(allTestClusters[0:2], nil)

		// check transporter
		vehicleRepo.EXPECT().GetByID(
			ctx,
			int32(2),
		).Return(nil, storage.ErrVehicleNotFound)

		// when
		result, err := svc.Update(ctx, wateringPlanID, updatedWateringPlan)

		// then
		assert.Nil(t, result)
		assert.EqualError(t, err, "404: vehicle not found")
	})

	t.Run("should return an error when finding treeclusters fails", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo)

		// check treecluster
		clusterRepo.EXPECT().GetByIDs(
			ctx,
			[]int32{1, 2},
		).Return(nil, storage.ErrConnectionClosed)

		// when
		result, err := svc.Update(ctx, wateringPlanID, updatedWateringPlan)

		// then
		assert.Nil(t, result)
		assert.EqualError(t, err, "500: connection is closed")
	})

	t.Run("should return an error when treecluster are empty", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo)

		// check treecluster
		clusterRepo.EXPECT().GetByIDs(
			ctx,
			[]int32{1, 2},
		).Return([]*entities.TreeCluster{}, nil)

		// when
		result, err := svc.Update(ctx, wateringPlanID, updatedWateringPlan)

		// then
		assert.Nil(t, result)
		assert.EqualError(t, err, "404: treecluster not found")
	})

	t.Run("should return an error when watering plan does not exist", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo)

		// check treecluster
		clusterRepo.EXPECT().GetByIDs(
			ctx,
			[]int32{1, 2},
		).Return(allTestClusters[0:2], nil)

		// check transporter
		vehicleRepo.EXPECT().GetByID(
			ctx,
			int32(2),
		).Return(allTestVehicles[1], nil)

		// check trailer
		vehicleRepo.EXPECT().GetByID(
			ctx,
			int32(1),
		).Return(allTestVehicles[0], nil)

		wateringPlanRepo.EXPECT().Update(
			ctx,
			wateringPlanID,
			mock.Anything,
		).Return(storage.ErrEntityNotFound)

		// when
		result, err := svc.Update(ctx, wateringPlanID, updatedWateringPlan)

		// then
		assert.Nil(t, result)
		assert.EqualError(t, err, "404: watering plan not found")
	})

	t.Run("should return an error when the update fails", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo)

		expectedErr := errors.New("failed to update watering plan")

		// check treecluster
		clusterRepo.EXPECT().GetByIDs(
			ctx,
			[]int32{1, 2},
		).Return(allTestClusters[0:2], nil)

		// check transporter
		vehicleRepo.EXPECT().GetByID(
			ctx,
			int32(2),
		).Return(allTestVehicles[1], nil)

		// check trailer
		vehicleRepo.EXPECT().GetByID(
			ctx,
			int32(1),
		).Return(allTestVehicles[0], nil)

		wateringPlanRepo.EXPECT().Update(
			ctx,
			wateringPlanID,
			mock.Anything,
		).Return(expectedErr)

		// when
		result, err := svc.Update(context.Background(), wateringPlanID, updatedWateringPlan)

		// then
		assert.Nil(t, result)
		assert.EqualError(t, err, "500: failed to update watering plan")
	})

	t.Run("should return validation error on empty date", func(t *testing.T) {
		// given
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo)

		updatedWateringPlan.Date = time.Time{}

		// when
		result, err := svc.Update(ctx, wateringPlanID, updatedWateringPlan)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validation error")
	})

	t.Run("should return validation error on empty transporter", func(t *testing.T) {
		// given
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo)

		updatedWateringPlan := &entities.WateringPlanUpdate{
			Date:        time.Date(2024, 8, 3, 0, 0, 0, 0, time.UTC),
			Description: "New watering plan for the east side of the city",
		}

		// when
		result, err := svc.Update(ctx, wateringPlanID, updatedWateringPlan)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validation error")
	})

	t.Run("should return validation error on empty treeclusters", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo)

		updatedWateringPlan := &entities.WateringPlanUpdate{
			Date:        time.Date(2024, 9, 26, 0, 0, 0, 0, time.UTC),
			Description: "Updated watering plan",
			TransporterID: utils.P(int32(2)),
		}

		// when
		result, err := svc.Update(ctx, wateringPlanID, updatedWateringPlan)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validation error")
	})
}

func TestWateringPlanService_Delete(t *testing.T) {
	ctx := context.Background()

	wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
	clusterRepo := storageMock.NewMockTreeClusterRepository(t)
	vehicleRepo := storageMock.NewMockVehicleRepository(t)
	svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo)

	t.Run("should successfully delete a watering plan", func(t *testing.T) {
		id := int32(1)

		wateringPlanRepo.EXPECT().GetByID(ctx, id).Return(allTestWateringPlans[1], nil)
		wateringPlanRepo.EXPECT().Delete(ctx, id).Return(nil)

		// when
		err := svc.Delete(ctx, id)

		// then
		assert.NoError(t, err)
	})

	t.Run("should return error if watering plan not found", func(t *testing.T) {
		id := int32(2)

		wateringPlanRepo.EXPECT().GetByID(ctx, id).Return(nil, storage.ErrEntityNotFound)

		// when
		err := svc.Delete(ctx, id)

		// then
		assert.Error(t, err)
		assert.EqualError(t, err, "404: watering plan not found")
	})

	t.Run("should return error if deleting watering plan fails", func(t *testing.T) {
		id := int32(4)

		wateringPlanRepo.EXPECT().GetByID(ctx, id).Return(allTestWateringPlans[1], nil)
		expectedErr := errors.New("failed to delete")
		wateringPlanRepo.EXPECT().Delete(ctx, id).Return(expectedErr)

		// when
		err := svc.Delete(ctx, id)

		// then
		assert.Error(t, err)
		assert.EqualError(t, err, "500: failed to delete")
	})
}

var allTestWateringPlans = []*entities.WateringPlan{
	{
		ID:                 1,
		Date:               time.Date(2024, 9, 22, 0, 0, 0, 0, time.UTC),
		Description:        "New watering plan for the west side of the city",
		Status:             entities.WateringPlanStatusPlanned,
		Distance:           utils.P(63.0),
		TotalWaterRequired: utils.P(6000.0),
		Transporter:        allTestVehicles[1],
		Trailer:            allTestVehicles[0],
		Treecluster:        allTestClusters[0:2],
	},
	{
		ID:                 2,
		Date:               time.Date(2024, 8, 3, 0, 0, 0, 0, time.UTC),
		Description:        "New watering plan for the east side of the city",
		Status:             entities.WateringPlanStatusActive,
		Distance:           utils.P(63.0),
		TotalWaterRequired: utils.P(6000.0),
		Transporter:        allTestVehicles[1],
		Trailer:            allTestVehicles[0],
		Treecluster:        allTestClusters[2:3],
	},
	{
		ID:                 3,
		Date:               time.Date(2024, 6, 12, 0, 0, 0, 0, time.UTC),
		Description:        "Very important watering plan due to no rainfall",
		Status:             entities.WateringPlanStatusFinished,
		Distance:           utils.P(63.0),
		TotalWaterRequired: utils.P(6000.0),
		Transporter:        allTestVehicles[1],
		Trailer:            nil,
		Treecluster:        allTestClusters[0:3],
	},
	{
		ID:                 4,
		Date:               time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC),
		Description:        "New watering plan for the south side of the city",
		Status:             entities.WateringPlanStatusNotCompeted,
		Distance:           utils.P(63.0),
		TotalWaterRequired: utils.P(6000.0),
		Transporter:        allTestVehicles[1],
		Trailer:            nil,
		Treecluster:        allTestClusters[2:3],
	},
	{
		ID:                 5,
		Date:               time.Date(2024, 6, 4, 0, 0, 0, 0, time.UTC),
		Description:        "Canceled due to flood",
		Status:             entities.WateringPlanStatusCanceled,
		Distance:           utils.P(63.0),
		TotalWaterRequired: utils.P(6000.0),
		Transporter:        allTestVehicles[1],
		Trailer:            nil,
		Treecluster:        allTestClusters[2:3],
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

var allTestClusters = []*entities.TreeCluster{
	{
		ID:             1,
		Name:           "Solitüde Strand",
		WateringStatus: entities.WateringStatusGood,
		MoistureLevel:  0.75,
		Region: &entities.Region{
			ID:   1,
			Name: "Mürwik",
		},
		Address:       "Solitüde Strand",
		Description:   "Alle Bäume am Strand",
		SoilCondition: entities.TreeSoilConditionSandig,
		Latitude:      utils.P(54.820940),
		Longitude:     utils.P(9.489022),
		Trees: []*entities.Tree{
			{ID: 1},
			{ID: 2},
			{ID: 3},
		},
	},
	{
		ID:             2,
		Name:           "Sankt-Jürgen-Platz",
		WateringStatus: entities.WateringStatusModerate,
		MoistureLevel:  0.5,
		Region: &entities.Region{
			ID:   1,
			Name: "Mürwik",
		},
		Address:       "Ulmenstraße",
		Description:   "Bäume beim Sankt-Jürgen-Platz",
		SoilCondition: entities.TreeSoilConditionSchluffig,
		Latitude:      utils.P(54.78805731048199),
		Longitude:     utils.P(9.44400186680097),
		Trees: []*entities.Tree{
			{ID: 4},
			{ID: 5},
			{ID: 6},
		},
	},
	{
		ID:             3,
		Name:           "Flensburger Stadion",
		WateringStatus: "unknown",
		MoistureLevel:  0.7,
		Region: &entities.Region{
			ID:   1,
			Name: "Mürwik",
		},
		Address:       "Flensburger Stadion",
		Description:   "Alle Bäume in der Gegend des Stadions in Mürwik",
		SoilCondition: "schluffig",
		Latitude:      utils.P(54.802163),
		Longitude:     utils.P(9.446398),
		Trees:         []*entities.Tree{},
	},
}

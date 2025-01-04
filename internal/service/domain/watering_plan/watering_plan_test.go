package wateringplan

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
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
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

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
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

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
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		expectedErr := errors.New("GetAll failed")
		wateringPlanRepo.EXPECT().GetAll(ctx).Return(nil, expectedErr)

		// when
		wateringPlans, err := svc.GetAll(ctx)

		// then
		assert.Error(t, err)
		assert.Nil(t, wateringPlans)
		assert.Equal(t, "500: GetAll failed (at internal/service/domain/watering_plan/watering_plan.go:210)", err.Error())
	})
}

func TestWateringPlanService_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("should return watering plan when found", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

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
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		id := int32(1)
		expectedErr := storage.ErrEntityNotFound
		wateringPlanRepo.EXPECT().GetByID(ctx, id).Return(nil, expectedErr)

		// when
		wateringPlan, err := svc.GetByID(ctx, id)

		// then
		assert.Error(t, err)
		assert.Nil(t, wateringPlan)
		assert.Equal(t, "404: watering plan not found (at internal/service/domain/watering_plan/watering_plan.go:207)", err.Error())
	})
}

func TestWateringPlanService_Create(t *testing.T) {
	ctx := context.Background()
	testUUIDString := "6a1078e8-80fd-458f-b74e-e388fe2dd6ab"
	testUUID, err := uuid.Parse(testUUIDString)
	if err != nil {
		t.Fatal(err)
	}

	newWateringPlan := &entities.WateringPlanCreate{
		Date:           time.Date(2024, 9, 26, 0, 0, 0, 0, time.UTC),
		Description:    "New watering plan",
		TransporterID:  utils.P(int32(2)),
		TrailerID:      utils.P(int32(1)),
		TreeClusterIDs: []*int32{utils.P(int32(1)), utils.P(int32(2))},
		UserIDs:        []*uuid.UUID{&testUUID},
	}

	t.Run("should successfully create a new watering plan", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		// check users
		userRepo.EXPECT().GetByIDs(
			ctx,
			[]string{testUUIDString},
		).Return([]*entities.User{testUser}, nil)

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
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		newWateringPlan := &entities.WateringPlanCreate{
			Date:           time.Date(2024, 9, 26, 0, 0, 0, 0, time.UTC),
			Description:    "New watering plan",
			TransporterID:  utils.P(int32(2)),
			TreeClusterIDs: []*int32{utils.P(int32(1)), utils.P(int32(2))},
			UserIDs:        []*uuid.UUID{&testUUID},
		}

		// check users
		userRepo.EXPECT().GetByIDs(
			ctx,
			[]string{testUUIDString},
		).Return([]*entities.User{testUser}, nil)

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

	t.Run("should return an error when finding users fails", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		// check treecluster
		userRepo.EXPECT().GetByIDs(
			ctx,
			[]string{testUUIDString},
		).Return(nil, storage.ErrUserNotFound)

		// when
		result, err := svc.Create(ctx, newWateringPlan)

		// then
		assert.Nil(t, result)
		assert.EqualError(t, err, "404: user not found")
	})

	t.Run("should return an error when finding treeclusters fails", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		// check users
		userRepo.EXPECT().GetByIDs(
			ctx,
			[]string{testUUIDString},
		).Return([]*entities.User{testUser}, nil)

		// check treecluster
		clusterRepo.EXPECT().GetByIDs(
			ctx,
			[]int32{1, 2},
		).Return(nil, storage.ErrConnectionClosed)

		// when
		result, err := svc.Create(ctx, newWateringPlan)

		// then
		assert.Nil(t, result)
		assert.EqualError(t, err, "500: connection is closed (at internal/service/domain/watering_plan/watering_plan.go:210)")
	})

	t.Run("should return an error when users are empty", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		// check users
		userRepo.EXPECT().GetByIDs(
			ctx,
			[]string{testUUIDString},
		).Return([]*entities.User{}, nil)

		// when
		result, err := svc.Create(ctx, newWateringPlan)

		// then
		assert.Nil(t, result)
		assert.EqualError(t, err, "404: user not found")
	})

	t.Run("should return an error when treecluster are empty", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		// check users
		userRepo.EXPECT().GetByIDs(
			ctx,
			[]string{testUUIDString},
		).Return([]*entities.User{testUser}, nil)

		// check treecluster
		clusterRepo.EXPECT().GetByIDs(
			ctx,
			[]int32{1, 2},
		).Return([]*entities.TreeCluster{}, nil)

		// when
		result, err := svc.Create(ctx, newWateringPlan)

		// then
		assert.Nil(t, result)
		assert.EqualError(t, err, "404: treecluster not found (at internal/service/domain/watering_plan/watering_plan.go:177)")
	})

	t.Run("should return an error when transporter is not found", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		// check users
		userRepo.EXPECT().GetByIDs(
			ctx,
			[]string{testUUIDString},
		).Return([]*entities.User{testUser}, nil)

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
		assert.EqualError(t, err, "404: vehicle not found (at internal/service/domain/watering_plan/watering_plan.go:165)")
	})

	t.Run("should return an error when creating watering plan fails", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		expectedErr := errors.New("Failed to create watering plan")

		// check users
		userRepo.EXPECT().GetByIDs(
			ctx,
			[]string{testUUIDString},
		).Return([]*entities.User{testUser}, nil)

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
		assert.EqualError(t, err, "500: Failed to create watering plan (at internal/service/domain/watering_plan/watering_plan.go:210)")
	})

	t.Run("should return validation error when TreeClusterIDs contains nil pointers", func(t *testing.T) {
		// given
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		newWateringPlan := &entities.WateringPlanCreate{
			Date:           time.Date(2024, 9, 26, 0, 0, 0, 0, time.UTC),
			Description:    "New watering plan with nil TreeClusterIDs",
			TransporterID:  utils.P(int32(2)),
			TreeClusterIDs: []*int32{nil, nil},
			UserIDs:        []*uuid.UUID{&testUUID},
		}

		// when
		result, err := svc.Create(ctx, newWateringPlan)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validation error")
	})

	t.Run("should return validation error on empty date", func(t *testing.T) {
		// given
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

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
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		newWateringPlan := &entities.WateringPlanCreate{
			Date:        time.Date(2024, 9, 26, 0, 0, 0, 0, time.UTC),
			Description: "New watering plan",
			UserIDs:     []*uuid.UUID{&testUUID},
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
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		newWateringPlan := &entities.WateringPlanCreate{
			Date:          time.Date(2024, 9, 26, 0, 0, 0, 0, time.UTC),
			Description:   "New watering plan",
			TransporterID: utils.P(int32(2)),
			UserIDs:       []*uuid.UUID{&testUUID},
		}

		// when
		result, err := svc.Create(ctx, newWateringPlan)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validation error")
	})

	t.Run("should return validation error on empty users", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		newWateringPlan := &entities.WateringPlanCreate{
			Date:          time.Date(2024, 9, 26, 0, 0, 0, 0, time.UTC),
			Description:   "New watering plan",
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
	testUUIDString := "6a1078e8-80fd-458f-b74e-e388fe2dd6ab"
	testUUID, err := uuid.Parse(testUUIDString)
	if err != nil {
		t.Fatal(err)
	}

	updatedWateringPlan := &entities.WateringPlanUpdate{
		Date:             time.Date(2024, 8, 3, 0, 0, 0, 0, time.UTC),
		Description:      "New watering plan for the east side of the city",
		TransporterID:    utils.P(int32(2)),
		TrailerID:        utils.P(int32(1)),
		TreeClusterIDs:   []*int32{utils.P(int32(1)), utils.P(int32(2))},
		UserIDs:          []*uuid.UUID{&testUUID},
		Status:           entities.WateringPlanStatusActive,
		CancellationNote: "",
	}

	t.Run("should successfully update a watering plan", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		// check users
		userRepo.EXPECT().GetByIDs(
			ctx,
			[]string{testUUIDString},
		).Return([]*entities.User{testUser}, nil)

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

	t.Run("should successfully update a watering plan with evaluation", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		updatedWateringPlan := &entities.WateringPlanUpdate{
			Date:             time.Date(2024, 8, 3, 0, 0, 0, 0, time.UTC),
			Status:           entities.WateringPlanStatusFinished,
			CancellationNote: "",
			Description:      "New watering plan for the east side of the city",
			TransporterID:    utils.P(int32(2)),
			TreeClusterIDs:   []*int32{utils.P(int32(1)), utils.P(int32(2))},
			UserIDs:          []*uuid.UUID{&testUUID},
			Evaluation: []*entities.EvaluationValue{
				{
					WateringPlanID: wateringPlanID,
					TreeClusterID:  1,
					ConsumedWater:  utils.P(100.00),
				},
			},
		}

		// check users
		userRepo.EXPECT().GetByIDs(
			ctx,
			[]string{testUUIDString},
		).Return([]*entities.User{testUser}, nil)

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
		).Return(allTestWateringPlans[2], nil)

		// when
		result, err := svc.Update(ctx, wateringPlanID, updatedWateringPlan)

		// then
		assert.NoError(t, err)
		assert.Equal(t, allTestWateringPlans[2], result)
	})

	t.Run("should successfully update a watering plan without a trailer", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		updatedWateringPlan := &entities.WateringPlanUpdate{
			Date:             time.Date(2024, 8, 3, 0, 0, 0, 0, time.UTC),
			Status:           entities.WateringPlanStatusActive,
			CancellationNote: "",
			Description:      "New watering plan for the east side of the city",
			TransporterID:    utils.P(int32(2)),
			TreeClusterIDs:   []*int32{utils.P(int32(1)), utils.P(int32(2))},
			UserIDs:          []*uuid.UUID{&testUUID},
		}

		// check users
		userRepo.EXPECT().GetByIDs(
			ctx,
			[]string{testUUIDString},
		).Return([]*entities.User{testUser}, nil)

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

	t.Run("should return an error when users is not found", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		// check users
		userRepo.EXPECT().GetByIDs(
			ctx,
			[]string{testUUIDString},
		).Return(nil, storage.ErrUserNotFound)

		// when
		result, err := svc.Update(ctx, wateringPlanID, updatedWateringPlan)

		// then
		assert.Nil(t, result)
		assert.EqualError(t, err, "404: user not found")
	})

	t.Run("should return an error when transporter is not found", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		// check users
		userRepo.EXPECT().GetByIDs(
			ctx,
			[]string{testUUIDString},
		).Return([]*entities.User{testUser}, nil)

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
		assert.EqualError(t, err, "404: vehicle not found (at internal/service/domain/watering_plan/watering_plan.go:165)")
	})

	t.Run("should return an error when finding treeclusters fails", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		// check users
		userRepo.EXPECT().GetByIDs(
			ctx,
			[]string{testUUIDString},
		).Return([]*entities.User{testUser}, nil)

		// check treecluster
		clusterRepo.EXPECT().GetByIDs(
			ctx,
			[]int32{1, 2},
		).Return(nil, storage.ErrConnectionClosed)

		// when
		result, err := svc.Update(ctx, wateringPlanID, updatedWateringPlan)

		// then
		assert.Nil(t, result)
		assert.EqualError(t, err, "500: connection is closed (at internal/service/domain/watering_plan/watering_plan.go:210)")
	})

	t.Run("should return an error when treecluster are empty", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		// check users
		userRepo.EXPECT().GetByIDs(
			ctx,
			[]string{testUUIDString},
		).Return([]*entities.User{testUser}, nil)

		// check treecluster
		clusterRepo.EXPECT().GetByIDs(
			ctx,
			[]int32{1, 2},
		).Return([]*entities.TreeCluster{}, nil)

		// when
		result, err := svc.Update(ctx, wateringPlanID, updatedWateringPlan)

		// then
		assert.Nil(t, result)
		assert.EqualError(t, err, "404: treecluster not found (at internal/service/domain/watering_plan/watering_plan.go:177)")
	})

	t.Run("should return an error when users are empty", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		// check users
		userRepo.EXPECT().GetByIDs(
			ctx,
			[]string{testUUIDString},
		).Return([]*entities.User{}, nil)

		// when
		result, err := svc.Update(ctx, wateringPlanID, updatedWateringPlan)

		// then
		assert.Nil(t, result)
		assert.EqualError(t, err, "404: user not found")
	})

	t.Run("should return an error when watering plan does not exist", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		// check users
		userRepo.EXPECT().GetByIDs(
			ctx,
			[]string{testUUIDString},
		).Return([]*entities.User{testUser}, nil)

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
		assert.EqualError(t, err, "404: watering plan not found (at internal/service/domain/watering_plan/watering_plan.go:207)")
	})

	t.Run("should return an error when the update fails", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		expectedErr := errors.New("failed to update watering plan")

		// check users
		userRepo.EXPECT().GetByIDs(
			ctx,
			[]string{testUUIDString},
		).Return([]*entities.User{testUser}, nil)

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
		assert.EqualError(t, err, "500: failed to update watering plan (at internal/service/domain/watering_plan/watering_plan.go:210)")
	})

	t.Run("should return error if cancellation note is not empty but the status is not »canceled«", func(t *testing.T) {
		// given
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		updatedWateringPlan.Status = entities.WateringPlanStatusActive
		updatedWateringPlan.CancellationNote = "This is a note"

		// when
		result, err := svc.Update(ctx, wateringPlanID, updatedWateringPlan)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "400: Cancellation note can only be set if watering plan is canceled (at internal/service/domain/watering_plan/watering_plan.go:195)")
	})

	t.Run("should return error if the evaluation is not empty but the status is not »finished«", func(t *testing.T) {
		// given
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		updatedWateringPlan.Status = entities.WateringPlanStatusPlanned
		updatedWateringPlan.CancellationNote = ""
		updatedWateringPlan.Evaluation = []*entities.EvaluationValue{
			{
				WateringPlanID: wateringPlanID,
				TreeClusterID:  1,
				ConsumedWater:  utils.P(100.00),
			},
		}

		// when
		result, err := svc.Update(ctx, wateringPlanID, updatedWateringPlan)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "400: Evaluation values can only be set if the watering plan has been finished (at internal/service/domain/watering_plan/watering_plan.go:199)")
	})

	t.Run("should return validation error when TreeClusterIDs contains nil pointers", func(t *testing.T) {
		// given
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		updatedWateringPlan := &entities.WateringPlanUpdate{
			Date:           time.Date(2024, 9, 26, 0, 0, 0, 0, time.UTC),
			Description:    "New watering plan with nil TreeClusterIDs",
			TransporterID:  utils.P(int32(2)),
			TreeClusterIDs: []*int32{nil, nil},
			UserIDs:        []*uuid.UUID{&testUUID},
		}

		// when
		result, err := svc.Update(ctx, wateringPlanID, updatedWateringPlan)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validation error")
	})

	t.Run("should return validation error when UserIDs contains nil pointers", func(t *testing.T) {
		// given
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		updatedWateringPlan := &entities.WateringPlanUpdate{
			Date:        time.Date(2024, 8, 3, 0, 0, 0, 0, time.UTC),
			Status:      entities.WateringPlanStatusActive,
			Description: "New watering plan for the east side of the city",
			UserIDs:     []*uuid.UUID{nil, nil},
		}

		// when
		result, err := svc.Update(ctx, wateringPlanID, updatedWateringPlan)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validation error")
	})

	t.Run("should return validation error on empty date", func(t *testing.T) {
		// given
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		updatedWateringPlan.Date = time.Time{}

		// when
		result, err := svc.Update(ctx, wateringPlanID, updatedWateringPlan)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validation error")
	})

	t.Run("should return validation error on wrong status", func(t *testing.T) {
		// given
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		updatedWateringPlan.Status = "test"

		// when
		result, err := svc.Update(ctx, wateringPlanID, updatedWateringPlan)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validation error")
	})

	t.Run("should return validation error on empty users", func(t *testing.T) {
		// given
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		updatedWateringPlan := &entities.WateringPlanUpdate{
			Date:        time.Date(2024, 8, 3, 0, 0, 0, 0, time.UTC),
			Status:      entities.WateringPlanStatusActive,
			Description: "New watering plan for the east side of the city",
			UserIDs:     []*uuid.UUID{&testUUID},
		}

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
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		updatedWateringPlan := &entities.WateringPlanUpdate{
			Date:        time.Date(2024, 8, 3, 0, 0, 0, 0, time.UTC),
			Status:      entities.WateringPlanStatusActive,
			Description: "New watering plan for the east side of the city",
			UserIDs:     []*uuid.UUID{&testUUID},
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
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

		updatedWateringPlan := &entities.WateringPlanUpdate{
			Date:          time.Date(2024, 9, 26, 0, 0, 0, 0, time.UTC),
			Description:   "Updated watering plan",
			Status:        entities.WateringPlanStatusActive,
			TransporterID: utils.P(int32(2)),
			UserIDs:       []*uuid.UUID{&testUUID},
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
	userRepo := storageMock.NewMockUserRepository(t)
	svc := NewWateringPlanService(wateringPlanRepo, clusterRepo, vehicleRepo, userRepo)

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
		assert.EqualError(t, err, "404: watering plan not found (at internal/service/domain/watering_plan/watering_plan.go:207)")
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
		assert.EqualError(t, err, "500: failed to delete (at internal/service/domain/watering_plan/watering_plan.go:210)")
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
		TreeClusters:       allTestClusters[0:2],
		CancellationNote:   "",
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
		TreeClusters:       allTestClusters[2:3],
		CancellationNote:   "",
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
		TreeClusters:       allTestClusters[0:3],
		CancellationNote:   "",
		Evaluation: []*entities.EvaluationValue{
			{
				WateringPlanID: 3,
				TreeClusterID:  1,
				ConsumedWater:  utils.P(10.0),
			},
			{
				WateringPlanID: 3,
				TreeClusterID:  2,
				ConsumedWater:  utils.P(10.0),
			},
			{
				WateringPlanID: 3,
				TreeClusterID:  3,
				ConsumedWater:  utils.P(10.0),
			},
		},
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
		TreeClusters:       allTestClusters[2:3],
		CancellationNote:   "",
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
		TreeClusters:       allTestClusters[2:3],
		CancellationNote:   "The watering plan was cancelled due to various reasons.",
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

var testUser = &entities.User{
	Username:    "user1",
	FirstName:   "John",
	LastName:    "Doe",
	Email:       "john.doe@green-ecolution.de",
	EmployeeID:  "EMP001",
	PhoneNumber: "+49 123456789",
}

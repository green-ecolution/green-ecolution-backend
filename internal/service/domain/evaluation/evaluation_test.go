package evaluation_test

import (
	"context"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service/domain/evaluation"
	storageMock "github.com/green-ecolution/green-ecolution-backend/internal/storage/_mock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestEvaluationService_GetAll(t *testing.T) {
	expectedEvaluation := &entities.Evaluation{
		TreeCount:             int64(10),
		TreeClusterCount:      int64(3),
		WateringPlanCount:     int64(3),
		SensorCount:           int64(2),
		TotalWaterConsumption: int64(10000),
	}

	t.Run("should return evaluation values when successful", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)

		svc := evaluation.NewEvaluationService(clusterRepo, treeRepo, sensorRepo, wateringPlanRepo)

		clusterRepo.EXPECT().GetCount(context.Background(), "").Return(expectedEvaluation.TreeClusterCount, nil)
		treeRepo.EXPECT().GetCount(context.Background(), "").Return(expectedEvaluation.TreeCount, nil)
		sensorRepo.EXPECT().GetCount(context.Background(), "").Return(expectedEvaluation.SensorCount, nil)
		wateringPlanRepo.EXPECT().GetCount(context.Background(), "").Return(expectedEvaluation.WateringPlanCount, nil)
		wateringPlanRepo.EXPECT().GetTotalConsumedWater(context.Background()).Return(expectedEvaluation.TotalWaterConsumption, nil)

		evaluation, err := svc.GetEvaluation(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expectedEvaluation, evaluation)
	})

	t.Run("should return error when getting cluster count fails", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)

		svc := evaluation.NewEvaluationService(clusterRepo, treeRepo, sensorRepo, wateringPlanRepo)

		clusterRepo.EXPECT().GetCount(context.Background(), "").Return(int64(0), errors.New("internal error"))
		evaluation, err := svc.GetEvaluation(context.Background())

		assert.Error(t, err)
		assert.Equal(t, int64(0), evaluation.SensorCount)
		assert.Equal(t, int64(0), evaluation.TreeClusterCount)
		assert.Equal(t, int64(0), evaluation.TreeCount)
		assert.Equal(t, int64(0), evaluation.WateringPlanCount)
		assert.Equal(t, int64(0), evaluation.TotalWaterConsumption)
	})

	t.Run("should return error when getting tree count fails", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)

		svc := evaluation.NewEvaluationService(clusterRepo, treeRepo, sensorRepo, wateringPlanRepo)

		clusterRepo.EXPECT().GetCount(context.Background(), "").Return(expectedEvaluation.TreeClusterCount, nil)
		treeRepo.EXPECT().GetCount(context.Background(), "").Return(int64(0), errors.New("internal error"))
		evaluation, err := svc.GetEvaluation(context.Background())

		assert.Error(t, err)
		assert.Equal(t, int64(0), evaluation.SensorCount)
		assert.Equal(t, int64(0), evaluation.TreeClusterCount)
		assert.Equal(t, int64(0), evaluation.TreeCount)
		assert.Equal(t, int64(0), evaluation.WateringPlanCount)
		assert.Equal(t, int64(0), evaluation.TotalWaterConsumption)
	})

	t.Run("should return error when getting sensor count fails", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)

		svc := evaluation.NewEvaluationService(clusterRepo, treeRepo, sensorRepo, wateringPlanRepo)

		clusterRepo.EXPECT().GetCount(context.Background(), "").Return(expectedEvaluation.TreeClusterCount, nil)
		treeRepo.EXPECT().GetCount(context.Background(), "").Return(expectedEvaluation.TreeCount, nil)
		sensorRepo.EXPECT().GetCount(context.Background(), "").Return(int64(0), errors.New("internal error"))
		evaluation, err := svc.GetEvaluation(context.Background())

		assert.Error(t, err)
		assert.Equal(t, int64(0), evaluation.SensorCount)
		assert.Equal(t, int64(0), evaluation.TreeClusterCount)
		assert.Equal(t, int64(0), evaluation.TreeCount)
		assert.Equal(t, int64(0), evaluation.WateringPlanCount)
		assert.Equal(t, int64(0), evaluation.TotalWaterConsumption)
	})

	t.Run("should return error when getting watering plans count fails", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)

		svc := evaluation.NewEvaluationService(clusterRepo, treeRepo, sensorRepo, wateringPlanRepo)

		clusterRepo.EXPECT().GetCount(context.Background(), "").Return(expectedEvaluation.TreeClusterCount, nil)
		treeRepo.EXPECT().GetCount(context.Background(), "").Return(expectedEvaluation.TreeCount, nil)
		sensorRepo.EXPECT().GetCount(context.Background(), "").Return(expectedEvaluation.SensorCount, nil)
		wateringPlanRepo.EXPECT().GetCount(context.Background(), "").Return(int64(0), errors.New("internal error"))

		evaluation, err := svc.GetEvaluation(context.Background())

		assert.Error(t, err)
		assert.Equal(t, int64(0), evaluation.SensorCount)
		assert.Equal(t, int64(0), evaluation.TreeClusterCount)
		assert.Equal(t, int64(0), evaluation.TreeCount)
		assert.Equal(t, int64(0), evaluation.WateringPlanCount)
		assert.Equal(t, int64(0), evaluation.TotalWaterConsumption)
	})

	t.Run("should return error when getting total water consumption fails", func(t *testing.T) {
		wateringPlanRepo := storageMock.NewMockWateringPlanRepository(t)
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)

		svc := evaluation.NewEvaluationService(clusterRepo, treeRepo, sensorRepo, wateringPlanRepo)

		clusterRepo.EXPECT().GetCount(context.Background(), "").Return(expectedEvaluation.TreeClusterCount, nil)
		treeRepo.EXPECT().GetCount(context.Background(), "").Return(expectedEvaluation.TreeCount, nil)
		sensorRepo.EXPECT().GetCount(context.Background(), "").Return(expectedEvaluation.SensorCount, nil)
		wateringPlanRepo.EXPECT().GetCount(context.Background(), "").Return(expectedEvaluation.WateringPlanCount, nil)
		wateringPlanRepo.EXPECT().GetTotalConsumedWater(context.Background()).Return(int64(0), errors.New("internal error"))

		evaluation, err := svc.GetEvaluation(context.Background())

		assert.Error(t, err)
		assert.Equal(t, int64(0), evaluation.SensorCount)
		assert.Equal(t, int64(0), evaluation.TreeClusterCount)
		assert.Equal(t, int64(0), evaluation.TreeCount)
		assert.Equal(t, int64(0), evaluation.WateringPlanCount)
		assert.Equal(t, int64(0), evaluation.TotalWaterConsumption)
	})
}

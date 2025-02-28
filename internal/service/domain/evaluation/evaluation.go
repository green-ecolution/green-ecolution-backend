package evaluation

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
)

type EvaluationService struct {
	treeClusterRepo  storage.TreeClusterRepository
	treeRepo         storage.TreeRepository
	sensorRepo       storage.SensorRepository
	wateringPlanRepo storage.WateringPlanRepository
}

func NewEvaluationService(
	treeClusterRepo storage.TreeClusterRepository,
	treeRepo storage.TreeRepository,
	sensorRepo storage.SensorRepository,
	wateringPlanRepo storage.WateringPlanRepository,
) service.EvaluationService {
	return &EvaluationService{
		treeClusterRepo:  treeClusterRepo,
		treeRepo:         treeRepo,
		sensorRepo:       sensorRepo,
		wateringPlanRepo: wateringPlanRepo,
	}
}

func (e *EvaluationService) GetAll(ctx context.Context) (*entities.Evaluation, error) {
	log := logger.GetLogger(ctx)
	return nil, nil
}

func (e *EvaluationService) Ready() bool {
	return e.treeClusterRepo != nil &&
		e.treeRepo != nil &&
		e.sensorRepo != nil &&
		e.wateringPlanRepo != nil
}

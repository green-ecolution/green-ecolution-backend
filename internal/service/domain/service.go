package domain

import (
	"github.com/green-ecolution/green-ecolution-backend/internal/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/service/domain/auth"
	"github.com/green-ecolution/green-ecolution-backend/internal/service/domain/evaluation"
	"github.com/green-ecolution/green-ecolution-backend/internal/service/domain/info"
	"github.com/green-ecolution/green-ecolution-backend/internal/service/domain/plugin"
	"github.com/green-ecolution/green-ecolution-backend/internal/service/domain/region"
	"github.com/green-ecolution/green-ecolution-backend/internal/service/domain/sensor"
	"github.com/green-ecolution/green-ecolution-backend/internal/service/domain/tree"
	"github.com/green-ecolution/green-ecolution-backend/internal/service/domain/treecluster"
	"github.com/green-ecolution/green-ecolution-backend/internal/service/domain/vehicle"
	wateringplan "github.com/green-ecolution/green-ecolution-backend/internal/service/domain/watering_plan"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/worker"
)

func NewService(cfg *config.Config, repos *storage.Repository, eventMananger *worker.EventManager) *service.Services {
	return &service.Services{
		InfoService:         info.NewInfoService(repos.Info),
		TreeService:         tree.NewTreeService(repos.Tree, repos.Sensor, repos.Image, repos.TreeCluster, eventMananger),
		AuthService:         auth.NewAuthService(repos.Auth, repos.User, &cfg.IdentityAuth),
		RegionService:       region.NewRegionService(repos.Region),
		TreeClusterService:  treecluster.NewTreeClusterService(repos.TreeCluster, repos.Tree, repos.Region, eventMananger),
		VehicleService:      vehicle.NewVehicleService(repos.Vehicle),
		SensorService:       sensor.NewSensorService(repos.Sensor, repos.Tree, eventMananger),
		PluginService:       plugin.NewPluginManager(repos.Auth),
		WateringPlanService: wateringplan.NewWateringPlanService(repos.WateringPlan, repos.TreeCluster, repos.Vehicle, repos.User, eventMananger, repos.Routing, repos.GpxBucket),
		EvaluationService:   evaluation.NewEvaluationService(repos.TreeCluster, repos.Tree, repos.Sensor, repos.WateringPlan),
	}
}

package domain

import (
	"github.com/green-ecolution/green-ecolution-backend/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/service/domain/auth"
	"github.com/green-ecolution/green-ecolution-backend/internal/service/domain/info"
	"github.com/green-ecolution/green-ecolution-backend/internal/service/domain/region"
	"github.com/green-ecolution/green-ecolution-backend/internal/service/domain/sensor"
	"github.com/green-ecolution/green-ecolution-backend/internal/service/domain/tree"
	"github.com/green-ecolution/green-ecolution-backend/internal/service/domain/treecluster"
	"github.com/green-ecolution/green-ecolution-backend/internal/service/domain/vehicle"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
)

func NewService(cfg *config.Config, repos *storage.Repository) *service.Services {
	geoLocator := treecluster.NewGeoLocation(repos.TreeCluster, repos.Tree, repos.Region)

	return &service.Services{
		InfoService:        info.NewInfoService(repos.Info),
		MqttService:        sensor.NewMqttService(repos.Sensor),
		TreeService:        tree.NewTreeService(repos.Tree, repos.Sensor, repos.Image, repos.TreeCluster, geoLocator),
		AuthService:        auth.NewAuthService(repos.Auth, repos.User, &cfg.IdentityAuth),
		RegionService:      region.NewRegionService(repos.Region),
		TreeClusterService: treecluster.NewTreeClusterService(repos.TreeCluster, repos.Tree, repos.Region, geoLocator),
		SensorService:      sensor.NewSensorService(repos.Sensor),
		VehicleService:     vehicle.NewVehicleService(repos.Vehicle),
	}
}

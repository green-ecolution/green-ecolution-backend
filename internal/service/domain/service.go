package domain

import (
	"github.com/green-ecolution/green-ecolution-backend/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/service/domain/info"
	"github.com/green-ecolution/green-ecolution-backend/internal/service/domain/sensor"
	"github.com/green-ecolution/green-ecolution-backend/internal/service/domain/tree"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
)

func NewService(cfg *config.Config, repositories *storage.Repository) *service.Services {
	return &service.Services{
		InfoService: info.NewInfoService(repositories.Info),
		MqttService: sensor.NewMqttService(repositories.Sensor),
		TreeService: tree.NewTreeService(repositories.Tree, repositories.Sensor),
	}
}

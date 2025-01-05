package openrouteservice

import (
	"github.com/green-ecolution/green-ecolution-backend/internal/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/pkg/ors"
)

func NewRepository(_ *config.Config) *storage.Repository {
	orsCfg := &ors.Configuration{
		Host:   "localhost:8080",
		Scheme: "http",
		Servers: ors.ServerConfigurations{
			{
				URL: "http://localhost:8080/ors",
			},
		},
		Debug: true,
	}

	routingRepo := NewRouteRepo(orsCfg)
	return &storage.Repository{
		Routing: routingRepo,
	}
}

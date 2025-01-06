package openrouteservice

import (
	"github.com/green-ecolution/green-ecolution-backend/internal/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
)

func NewRepository(_ *config.Config) *storage.Repository {
	routingRepo := NewRouteRepo()
	return &storage.Repository{
		Routing: routingRepo,
	}
}

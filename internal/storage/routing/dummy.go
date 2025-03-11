package routing

import (
	"context"
	"io"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
)

type DummyRoutingRepo struct{}

func NewDummyRoutingRepo() *DummyRoutingRepo {
	return &DummyRoutingRepo{}
}

func (r *DummyRoutingRepo) GenerateRoute(ctx context.Context, vehicle *entities.Vehicle, clusters []*entities.TreeCluster) (*entities.GeoJSON, error) {
	return nil, storage.ErrRoutingServiceDisabled
}

func (r *DummyRoutingRepo) GenerateRawGpxRoute(ctx context.Context, vehicle *entities.Vehicle, clusters []*entities.TreeCluster) (io.ReadCloser, error) {
	return nil, storage.ErrRoutingServiceDisabled
}

func (r *DummyRoutingRepo) GenerateRouteInformation(ctx context.Context, vehicle *entities.Vehicle, clusters []*entities.TreeCluster) (*entities.RouteMetadata, error) {
	return nil, storage.ErrRoutingServiceDisabled
}

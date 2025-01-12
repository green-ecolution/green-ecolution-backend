package valhalla

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/url"

	"github.com/green-ecolution/green-ecolution-backend/internal/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/routing"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/routing/valhalla/valhalla"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/routing/vroom"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
)

// validate is RouteRepo implements storage.RoutingRepository
var _ storage.RoutingRepository = (*RouteRepo)(nil)

type RouteRepoConfig struct {
	routing config.RoutingConfig
}

type RouteRepo struct {
	vroom    vroom.VroomClient
	valhalla valhalla.ValhallaClient
	cfg      *RouteRepoConfig
}

func NewRouteRepo(cfg *RouteRepoConfig) (*RouteRepo, error) {
	vroomURL, err := url.Parse(cfg.routing.Valhalla.Optimization.Vroom.Host)
	if err != nil {
		return nil, err
	}
	valhallaURL, err := url.Parse(cfg.routing.Valhalla.Host)
	if err != nil {
		return nil, err
	}

	vroomClient := vroom.NewVroomClient(
		vroom.WithHostURL(vroomURL),
		vroom.WithStartPoint(cfg.routing.StartPoint),
		vroom.WithEndPoint(cfg.routing.EndPoint),
		vroom.WithWateringPoint(cfg.routing.WateringPoint),
	)
	valhalllaClient := valhalla.NewValhallaClient(
		valhalla.WithHostURL(valhallaURL),
	)

	return &RouteRepo{
		vroom:    vroomClient,
		valhalla: valhalllaClient,
		cfg:      cfg,
	}, nil
}

func (r *RouteRepo) GenerateRoute(ctx context.Context, vehicle *entities.Vehicle, clusters []*entities.TreeCluster) (*entities.GeoJSON, error) {
	_, route, err := r.prepareRoute(ctx, vehicle, clusters)
	if err != nil {
		return nil, err
	}

	entity, err := r.valhalla.DirectionsGeoJSON(ctx, route)
	if err != nil {
		return nil, err
	}

	metadata, err := routing.ConvertLocations(&r.cfg.routing)
	if err != nil {
		return nil, err
	}

	entity.Metadata = *metadata

	return entity, nil
}

func (r *RouteRepo) GenerateRawGpxRoute(ctx context.Context, vehicle *entities.Vehicle, clusters []*entities.TreeCluster) (io.ReadCloser, error) {
	_, route, err := r.prepareRoute(ctx, vehicle, clusters)
	if err != nil {
		return nil, err
	}

	return r.valhalla.DirectionsRawGpx(ctx, route)
}

func (r *RouteRepo) GenerateRouteInformation(ctx context.Context, vehicle *entities.Vehicle, clusters []*entities.TreeCluster) (*entities.RouteMetadata, error) {
	optimizedRoutes, route, err := r.prepareRoute(ctx, vehicle, clusters)
	if err != nil {
		return nil, err
	}

	// currently handle only the first route
	var refillCount int
	if len(optimizedRoutes.Routes) > 0 {
		oRoute := optimizedRoutes.Routes[0]
		reducedSteps := utils.Reduce(oRoute.Steps, r.reduceSteps, make([]*vroom.VroomRouteStep, 0, len(oRoute.Steps)))
		refillCount = len(utils.Filter(reducedSteps, func(step *vroom.VroomRouteStep) bool {
			return step.Type == "pickup"
		}))
	}

	rawDirections, err := r.valhalla.DirectionsJSON(ctx, route)
	if err != nil {
		return nil, err
	}

	return &entities.RouteMetadata{
		Refills:  refillCount,
		Distance: rawDirections.Trip.Summary.Length,
		Time:     rawDirections.Trip.Summary.Time,
	}, nil
}

func (r *RouteRepo) prepareRoute(ctx context.Context, vehicle *entities.Vehicle, clusters []*entities.TreeCluster) (optimized *vroom.VroomResponse, routes *valhalla.DirectionRequest, err error) {
	optimizedRoutes, err := r.vroom.OptimizeRoute(ctx, vehicle, clusters)
	if err != nil {
		slog.Error("failed to optimize route", "error", err)
		return nil, nil, err
	}

	// currently handle only the first route
	if len(optimizedRoutes.Routes) == 0 {
		slog.Error("there are no routes in vroom response", "routes", optimizedRoutes)
		return nil, nil, errors.New("empty routes")
	}
	oRoute := optimizedRoutes.Routes[0]
	reducedSteps := utils.Reduce(oRoute.Steps, r.reduceSteps, make([]*vroom.VroomRouteStep, 0, len(oRoute.Steps)))
	locations := utils.Map(reducedSteps, func(step *vroom.VroomRouteStep) valhalla.Location {
		return valhalla.Location{
			Lat:  step.Location[1],
			Lon:  step.Location[0],
			Type: "break",
		}
	})

	costingOpts := make(map[string]valhalla.CostingOptions)
	costingOpts["truck"] = valhalla.CostingOptions{
		Width:     vehicle.Width,
		Height:    vehicle.Height,
		Length:    vehicle.Length,
		Weight:    vehicle.Weight,
		AxleLoad:  0.0,
		AxleCount: 2,
	}

	return optimizedRoutes, &valhalla.DirectionRequest{
		Locations:      locations,
		Costing:        "truck",
		CostingOptions: costingOpts,
	}, nil
}

// Reduce multiple pickups to one
// "start" -> "pickup" -> "pickup" -> "delivery" => "start" -> "pickup" -> "delivery"
func (r *RouteRepo) reduceSteps(acc []*vroom.VroomRouteStep, current vroom.VroomRouteStep) []*vroom.VroomRouteStep {
	if len(acc) == 0 {
		return append(acc, &current)
	}

	prev := acc[len(acc)-1]
	if prev.Type != "pickup" {
		return append(acc, &current)
	}

	if current.Type != "pickup" {
		return append(acc, &current)
	}

	prev.Load = current.Load
	return acc
}

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
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/routing/openrouteservice/vroom"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/routing/valhalla/valhalla"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
)

// validate is RouteRepo implements storage.RoutingRepository
var _ storage.RoutingRepository = (*RouteRepo)(nil)

const (
	treeAmount int32 = 40 // how much water does a cluster need
)

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
	orsRoute, err := r.prepareRoute(ctx, vehicle, clusters)
	if err != nil {
		return nil, err
	}

	return r.valhalla.DirectionsGeoJSON(ctx, orsRoute)
}

func (r *RouteRepo) GenerateRawGpxRoute(ctx context.Context, vehicle *entities.Vehicle, clusters []*entities.TreeCluster) (io.ReadCloser, error) {
	route, err := r.prepareRoute(ctx, vehicle, clusters)
	if err != nil {
		return nil, err
	}

	return r.valhalla.DirectionsRawGpx(ctx, route)
}

func (r *RouteRepo) prepareRoute(ctx context.Context, vehicle *entities.Vehicle, clusters []*entities.TreeCluster) (*valhalla.DirectionRequest, error) {
	optimizedRoutes, err := r.optimizeRoute(ctx, vehicle, clusters)
	if err != nil {
		slog.Error("failed to optimize route", "error", err)
		return nil, err
	}

	// currently handle only the first route
	if len(optimizedRoutes.Routes) == 0 {
		slog.Error("there are no routes in vroom response", "routes", optimizedRoutes)
		return nil, errors.New("empty routes")
	}
	oRoute := optimizedRoutes.Routes[0]

	// Reduce multiple pickups to one
	// "start" -> "pickup" -> "pickup" -> "delivery" => "start" -> "pickup" -> "delivery"
	reducedSteps := utils.Reduce(oRoute.Steps, func(acc []*vroom.VroomRouteStep, current vroom.VroomRouteStep) []*vroom.VroomRouteStep {
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
	}, make([]*vroom.VroomRouteStep, 0, len(oRoute.Steps)))

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
		Weight:    3.5, // TODO: use real value
		AxleLoad:  0.0,
		AxleCount: 2,
	}

	return &valhalla.DirectionRequest{
		Locations:      locations,
		Costing:        "truck",
		CostingOptions: costingOpts,
	}, nil
}

func (r *RouteRepo) optimizeRoute(ctx context.Context, vehicle *entities.Vehicle, cluster []*entities.TreeCluster) (*vroom.VroomResponse, error) {
	vroomVehicle, err := r.toVroomVehicle(vehicle)
	if err != nil {
		if errors.Is(err, storage.ErrUnknownVehicleType) {
			slog.Error("unknown vehicle type. please specify vehicle type", "error", err, "vehicle_type", vehicle.Type)
		}

		return nil, err
	}

	shipments := r.toVroomShipments(cluster)
	req := &vroom.VroomReq{
		Vehicles:  []vroom.VroomVehicle{*vroomVehicle},
		Shipments: shipments,
	}

	resp, err := r.vroom.Send(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *RouteRepo) toVroomShipments(cluster []*entities.TreeCluster) []vroom.VroomShipments {
	// ignore tree cluster with empty coordinates
	filteredClusters := utils.Filter(cluster, func(c *entities.TreeCluster) bool {
		return c.Longitude != nil && c.Latitude != nil
	})

	nextID := int32(0)
	return utils.Map(filteredClusters, func(c *entities.TreeCluster) vroom.VroomShipments {
		shipment := vroom.VroomShipments{
			Amount: []int32{treeAmount},
			Pickup: vroom.VroomShipmentStep{
				ID:       nextID,
				Location: r.cfg.routing.WateringPoint,
			},
			Delivery: vroom.VroomShipmentStep{
				Description: c.Name,
				ID:          nextID + 1,
				Location:    []float64{*c.Longitude, *c.Latitude},
			},
		}

		nextID += 2
		return shipment
	})
}

func (r *RouteRepo) toVroomVehicle(vehicle *entities.Vehicle) (*vroom.VroomVehicle, error) {
	vehicleType, err := r.toOrsVehicleType(vehicle.Type)
	if err != nil {
		return nil, err
	}

	return &vroom.VroomVehicle{
		ID:          vehicle.ID,
		Description: vehicle.Description,
		Profile:     vehicleType,
		Start:       r.cfg.routing.StartPoint,
		End:         r.cfg.routing.EndPoint,
		Capacity:    []int32{int32(vehicle.WaterCapacity)}, // vroom don't accept floats
	}, nil
}

func (r *RouteRepo) toOrsVehicleType(_ entities.VehicleType) (string, error) {
	return "driving-car", nil
	// switch vecType {
	// case entities.VehicleTypeTrailer:
	// 	return "driving-car", nil

	// case entities.VehicleTypeTransporter:
	// 	return "driving-hgv", nil

	// default:
	// 	return "", storage.ErrUnknownVehicleType
	// }
}

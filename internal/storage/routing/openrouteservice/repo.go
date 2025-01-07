package openrouteservice

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/url"

	"github.com/green-ecolution/green-ecolution-backend/internal/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/routing/openrouteservice/ors"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/routing/openrouteservice/vroom"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
)

// validate is RouteRepo implements storage.RoutingRepository
var _ storage.RoutingRepository = (*RouteRepo)(nil)

const (
	treeAmount int32 = 40 // how much water does a cluster need
)

type RouteRepoConfig struct {
	routing config.RoutingConfig
	s3      config.S3Config
}

type RouteRepo struct {
	vroom vroom.VroomClient
	ors   ors.OrsClient
	cfg   *RouteRepoConfig
}

func NewRouteRepo(cfg *RouteRepoConfig) (*RouteRepo, error) {
	vroomURL, err := url.Parse(cfg.routing.Ors.Optimization.Vroom.Host)
	if err != nil {
		return nil, err
	}
	orsURL, err := url.Parse(cfg.routing.Ors.Host)
	if err != nil {
		return nil, err
	}

	vroomClient := vroom.NewVroomClient(
		vroom.WithHostURL(vroomURL),
	)
	orsClient := ors.NewOrsClient(
		ors.WithHostURL(orsURL),
	)

	return &RouteRepo{
		vroom: vroomClient,
		ors:   orsClient,
		cfg:   cfg,
	}, nil
}

func (r *RouteRepo) GenerateRoute(ctx context.Context, vehicle *entities.Vehicle, clusters []*entities.TreeCluster) (*entities.GeoJSON, error) {
	orsProfile, err := r.toOrsVehicleType(vehicle.Type)
	if err != nil {
		return nil, err
	}

	orsRoute, err := r.prepareOrsRoute(ctx, vehicle, clusters)
	if err != nil {
		return nil, err
	}

	return r.ors.DirectionsGeoJSON(ctx, orsProfile, orsRoute)
}

func (r *RouteRepo) GenerateRawGpxRoute(ctx context.Context, vehicle *entities.Vehicle, clusters []*entities.TreeCluster) (io.ReadCloser, error) {
	orsProfile, err := r.toOrsVehicleType(vehicle.Type)
	if err != nil {
		return nil, err
	}

	orsRoute, err := r.prepareOrsRoute(ctx, vehicle, clusters)
	if err != nil {
		return nil, err
	}

	return r.ors.DirectionsRawGpx(ctx, orsProfile, orsRoute)
}

func (r *RouteRepo) prepareOrsRoute(ctx context.Context, vehicle *entities.Vehicle, clusters []*entities.TreeCluster) (*ors.OrsDirectionRequest, error) {
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

	orsLocation := utils.Reduce(reducedSteps, func(acc [][]float64, current *vroom.VroomRouteStep) [][]float64 {
		return append(acc, current.Location)
	}, make([][]float64, 0, len(reducedSteps)))

	return &ors.OrsDirectionRequest{
		Coordinates: orsLocation,
		Units:       "m",
		Language:    "de-de",
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

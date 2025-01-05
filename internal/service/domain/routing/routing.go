package routing

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
)

var _ service.RoutingService = (*RoutingService)(nil)

type RoutingService struct {
	routingRepo storage.RoutingRepository
	clusterRepo storage.TreeClusterRepository
	vehicleRepo storage.VehicleRepository
}

func NewRoutingService(routingRepo storage.RoutingRepository, clusterRepo storage.TreeClusterRepository, vehicleRepo storage.VehicleRepository) *RoutingService {
	return &RoutingService{
		routingRepo: routingRepo,
		clusterRepo: clusterRepo,
		vehicleRepo: vehicleRepo,
	}
}

func (s *RoutingService) PreviewRoute(ctx context.Context, vehicleID int32, clusterIDs []int32) (*entities.GeoJSON, error) {
	vehicle, err := s.vehicleRepo.GetByID(ctx, vehicleID)
	if err != nil {
		slog.Error("can't find vehicle to preview route", "error", err, "vehicle_id", vehicleID)
		return nil, service.NewError(service.NotFound, fmt.Sprintf("vehicle with id %d not found", vehicleID))
	}

	clusters, err := s.clusterRepo.GetByIDs(ctx, clusterIDs)
	if err != nil {
		// when error, something is wrong with the db, else clusters should be an empty array
		return nil, err
	}

	geoJson, err := s.routingRepo.GenerateRoute(ctx, vehicle, clusters)
	if err != nil {
		if errors.Is(err, storage.ErrUnknownVehicleType) {
			slog.Error("the vehicle type is not supported", "error", err, "vehicle_type", vehicle.Type)
			return nil, service.NewError(service.NotFound, "vehicle type is not supported")
		}
		return nil, err
	}

	return geoJson, nil
}

func (s *RoutingService) Ready() bool {
	// TODO: check if ors service is ready/healthy
	return true
}

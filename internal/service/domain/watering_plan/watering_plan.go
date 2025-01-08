package wateringplan

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/google/uuid"

	"github.com/go-playground/validator/v10"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/green-ecolution/green-ecolution-backend/internal/worker"
)

type WateringPlanService struct {
	wateringPlanRepo storage.WateringPlanRepository
	clusterRepo      storage.TreeClusterRepository
	vehicleRepo      storage.VehicleRepository
	userRepo         storage.UserRepository
	routingRepo      storage.RoutingRepository
	gpxBucket        storage.S3Repository
	validator        *validator.Validate
	eventManager     *worker.EventManager
}

func NewWateringPlanService(
	wateringPlanRepository storage.WateringPlanRepository,
	clusterRepository storage.TreeClusterRepository,
	vehicleRepository storage.VehicleRepository,
	userRepository storage.UserRepository,
	eventManager *worker.EventManager,
	routingRepo storage.RoutingRepository,
	gpxRepo storage.S3Repository,
) service.WateringPlanService {
	return &WateringPlanService{
		wateringPlanRepo: wateringPlanRepository,
		clusterRepo:      clusterRepository,
		vehicleRepo:      vehicleRepository,
		userRepo:         userRepository,
		routingRepo:      routingRepo,
		gpxBucket:        gpxRepo,
		validator:        validator.New(),
		eventManager:     eventManager,
	}
}

func (w *WateringPlanService) publishUpdateEvent(ctx context.Context, prevWp *entities.WateringPlan) error {
	slog.Debug("publish new event", "event", entities.EventTypeUpdateWateringPlan, "service", "WateringPlanService")
	updatedWp, err := w.GetByID(ctx, prevWp.ID)
	if err != nil {
		return err
	}
	event := entities.NewEventUpdateWateringPlan(prevWp, updatedWp)
	if err := w.eventManager.Publish(event); err != nil {
		slog.Error("error while sending event after updating watering plan", "err", err, "watering_plan_id", prevWp.ID)
	}

	return nil
}

func (w *WateringPlanService) PreviewRoute(ctx context.Context, transporterID int32, trailerID *int32, clusterIDs []int32) (*entities.GeoJSON, error) {
	transporter, err := w.vehicleRepo.GetByID(ctx, transporterID)
	if err != nil {
		slog.Error("can't find selected transporter to preview route", "error", err, "vehicle_id", transporterID)
		return nil, service.NewError(service.NotFound, fmt.Sprintf("vehicle with id %d not found", transporterID))
	}

	var trailer *entities.Vehicle
	if trailerID != nil {
		trailer, err = w.vehicleRepo.GetByID(ctx, *trailerID)
		if err != nil {
			slog.Error("can't find selected trailer to preview route. route will be calculated without trailer", "error", err, "trailer_id", trailerID)
		}
	}

	clusters, err := w.clusterRepo.GetByIDs(ctx, clusterIDs)
	if err != nil {
		// when error, something is wrong with the db, else clusters should be an empty array
		return nil, err
	}

	geoJSON, err := w.routingRepo.GenerateRoute(ctx, w.mergeVehicle(transporter, trailer), clusters)
	if err != nil {
		if errors.Is(err, storage.ErrUnknownVehicleType) {
			slog.Error("the vehicle type is not supported", "error", err, "vehicle_type", transporter.Type)
			return nil, service.NewError(service.NotFound, "vehicle type is not supported")
		}
		return nil, err
	}

	return geoJSON, nil
}

func (w *WateringPlanService) GetAll(ctx context.Context) ([]*entities.WateringPlan, error) {
	plans, err := w.wateringPlanRepo.GetAll(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	return plans, nil
}

func (w *WateringPlanService) GetByID(ctx context.Context, id int32) (*entities.WateringPlan, error) {
	got, err := w.wateringPlanRepo.GetByID(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}

	return got, nil
}

func (w *WateringPlanService) Create(ctx context.Context, createWp *entities.WateringPlanCreate) (*entities.WateringPlan, error) {
	if err := w.validator.Struct(createWp); err != nil {
		return nil, service.NewError(service.BadRequest, errors.Join(err, errors.New("validation error")).Error())
	}

	// TODO: get distance from valhalla
	// TODO: validate driver license

	if err := w.validateUserIDs(ctx, createWp.UserIDs); err != nil {
		return nil, service.NewError(service.NotFound, storage.ErrUserNotFound.Error())
	}

	treeClusters, err := w.fetchTreeClusters(ctx, createWp.TreeClusterIDs)
	if err != nil {
		return nil, err
	}

	transporter, err := w.fetchVehicle(ctx, *createWp.TransporterID)
	if err != nil {
		return nil, err
	}

	var trailer *entities.Vehicle
	if createWp.TrailerID != nil {
		trailer, err = w.fetchVehicle(ctx, *createWp.TrailerID)
		if err != nil {
			return nil, err // Maybe ignore?
		}
	}

	neededWater := w.calculateRequiredWater(treeClusters)
	created, err := w.wateringPlanRepo.Create(ctx, func(wp *entities.WateringPlan) (bool, error) {
		wp.Date = createWp.Date
		wp.Description = createWp.Description
		wp.Transporter = transporter
		wp.Trailer = trailer
		wp.TreeClusters = treeClusters
		wp.UserIDs = createWp.UserIDs
		wp.TotalWaterRequired = utils.P(float64(neededWater))

		return true, nil
	})
	if err != nil {
		return nil, handleError(err)
	}

	err = w.wateringPlanRepo.Update(ctx, created.ID, func(wp *entities.WateringPlan) (bool, error) {
		gpxURL, err := w.getGpxRouteURL(ctx, created.ID, w.mergeVehicle(transporter, trailer), treeClusters)
		if err != nil {
			return false, handleError(err)
		}

		wp.GpxURL = gpxURL

		return true, nil
	})

	if err != nil {
		return nil, handleError(err)
	}

	return created, nil
}

func (w *WateringPlanService) getGpxRouteURL(ctx context.Context, waterPlanID int32, vehicle *entities.Vehicle, clusters []*entities.TreeCluster) (string, error) {
	r, err := w.routingRepo.GenerateRawGpxRoute(ctx, vehicle, clusters)
	if err != nil {
		return "", err
	}
	defer r.Close()

	objName := fmt.Sprintf("waterplan-%d.gpx", waterPlanID)

	if err := w.gpxBucket.PutObject(ctx, objName, "application/gpx+xml;charset=UTF-8 ", -1, r); err != nil {
		return "", err
	}

	return fmt.Sprintf("/v1/watering-plan/route/gpx/%s", objName), nil
}

func (w *WateringPlanService) GetGPXFileStream(ctx context.Context, objName string) (io.ReadSeekCloser, error) {
	return w.gpxBucket.GetObject(ctx, objName)
}

func (w *WateringPlanService) Update(ctx context.Context, id int32, updateWp *entities.WateringPlanUpdate) (*entities.WateringPlan, error) {
	if err := w.validator.Struct(updateWp); err != nil {
		return nil, service.NewError(service.BadRequest, errors.Join(err, errors.New("validation error")).Error())
	}

	prevWp, err := w.GetByID(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}

	// TODO: get distance from valhalla
	// TODO: validate driver license

	if err := w.validateStatusDependentValues(updateWp); err != nil {
		return nil, err
	}

	if err := w.validateUserIDs(ctx, updateWp.UserIDs); err != nil {
		return nil, service.NewError(service.NotFound, storage.ErrUserNotFound.Error())
	}

	treeClusters, err := w.fetchTreeClusters(ctx, updateWp.TreeClusterIDs)
	if err != nil {
		return nil, err
	}

	transporter, err := w.fetchVehicle(ctx, *updateWp.TransporterID)
	if err != nil {
		return nil, err
	}

	var trailer *entities.Vehicle
	if updateWp.TrailerID != nil {
		trailer, err = w.fetchVehicle(ctx, *updateWp.TrailerID)
		if err != nil {
			return nil, err
		}
	}

	neededWater := w.calculateRequiredWater(treeClusters)
	err = w.wateringPlanRepo.Update(ctx, id, func(wp *entities.WateringPlan) (bool, error) {
		wp.Date = updateWp.Date
		wp.Description = updateWp.Description
		wp.Transporter = transporter
		wp.Trailer = trailer
		wp.TreeClusters = treeClusters
		wp.Status = updateWp.Status
		wp.CancellationNote = updateWp.CancellationNote
		wp.Evaluation = updateWp.Evaluation
		wp.UserIDs = updateWp.UserIDs
		wp.TotalWaterRequired = &neededWater

		return true, nil
	})

	if err != nil {
		return nil, handleError(err)
	}

	if err := w.publishUpdateEvent(ctx, prevWp); err != nil {
		return nil, handleError(err)
	}

	return w.GetByID(ctx, id)
}

func (w *WateringPlanService) Delete(ctx context.Context, id int32) error {
	_, err := w.wateringPlanRepo.GetByID(ctx, id)
	if err != nil {
		return handleError(err)
	}

	if err := w.wateringPlanRepo.Delete(ctx, id); err != nil {
		return handleError(err)
	}

	return nil
}

func (w *WateringPlanService) Ready() bool {
	return w.wateringPlanRepo != nil
}

func (w *WateringPlanService) fetchVehicle(ctx context.Context, vehicleID int32) (*entities.Vehicle, error) {
	vehicle, err := w.vehicleRepo.GetByID(ctx, vehicleID)
	if err != nil {
		return nil, service.NewError(service.NotFound, storage.ErrVehicleNotFound.Error())
	}

	return vehicle, nil
}

func (w *WateringPlanService) fetchTreeClusters(ctx context.Context, treeClusterIDs []*int32) ([]*entities.TreeCluster, error) {
	clusters, err := w.getTreeClusters(ctx, treeClusterIDs)
	if err != nil {
		return nil, handleError(err)
	}
	if len(clusters) == 0 {
		return nil, service.NewError(service.NotFound, storage.ErrTreeClusterNotFound.Error())
	}

	return clusters, nil
}

func (w *WateringPlanService) getTreeClusters(ctx context.Context, ids []*int32) ([]*entities.TreeCluster, error) {
	clusterIDs := make([]int32, len(ids))
	for i, id := range ids {
		clusterIDs[i] = *id
	}

	return w.clusterRepo.GetByIDs(ctx, clusterIDs)
}

// Checks if the incoming user ids are belonging to valid users
func (w *WateringPlanService) validateUserIDs(ctx context.Context, userIDs []*uuid.UUID) error {
	var userIDStrings []string
	for _, id := range userIDs {
		if id != nil {
			userIDStrings = append(userIDStrings, utils.UUIDToString(*id))
		}
	}

	users, err := w.userRepo.GetByIDs(ctx, userIDStrings)
	if err != nil {
		return handleError(err)
	}

	if len(users) == 0 {
		return storage.ErrUserNotFound
	}

	return nil
}

func (w *WateringPlanService) validateStatusDependentValues(entity *entities.WateringPlanUpdate) error {
	// Set cancellation note to nothing if the current status is not fitting
	if entity.CancellationNote != "" && entity.Status != entities.WateringPlanStatusCanceled {
		return service.NewError(service.BadRequest, "Cancellation note can only be set if watering plan is canceled")
	}

	if entity.Status != entities.WateringPlanStatusFinished && len(entity.Evaluation) > 0 {
		return service.NewError(service.BadRequest, "Evaluation values can only be set if the watering plan has been finished")
	}

	return nil
}

// This function calculates approximately how much water the irrigation schedule needs
// Each tree in a linked tree cluster requires approximately 120 liters of water
func (w *WateringPlanService) calculateRequiredWater(clusters []*entities.TreeCluster) float64 {
	return utils.Reduce(clusters, func(acc float64, tc *entities.TreeCluster) float64 {
		return acc + (float64(len(tc.Trees)) * 120.0)
	}, 0)
}

func (w *WateringPlanService) mergeVehicle(transporter, trailer *entities.Vehicle) *entities.Vehicle {
	if transporter == nil {
		return nil // this should not happen because of before validation
	}

	if trailer == nil {
		return transporter
	}

	return &entities.Vehicle{
		Width:         transporter.Width + trailer.Width,
		Height:        transporter.Height + trailer.Height,
		Length:        transporter.Length + trailer.Length,
		WaterCapacity: transporter.WaterCapacity + trailer.WaterCapacity,
		Type:          entities.VehicleTypeTransporter,
		NumberPlate:   fmt.Sprintf("%s - %s", transporter.NumberPlate, trailer.NumberPlate),
	}
}

func handleError(err error) error {
	if errors.Is(err, storage.ErrEntityNotFound) {
		return service.NewError(service.NotFound, storage.ErrWateringPlanNotFound.Error())
	}

	return service.NewError(service.InternalError, err.Error())
}

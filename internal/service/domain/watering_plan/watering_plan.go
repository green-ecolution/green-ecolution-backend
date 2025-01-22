package wateringplan

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/google/uuid"
	"github.com/spf13/viper"

	"github.com/go-playground/validator/v10"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
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
	log := logger.GetLogger(ctx)
	log.Debug("publish new event", "event", entities.EventTypeUpdateWateringPlan, "service", "WateringPlanService")
	updatedWp, err := w.GetByID(ctx, prevWp.ID)
	if err != nil {
		return err
	}
	event := entities.NewEventUpdateWateringPlan(prevWp, updatedWp)
	if err := w.eventManager.Publish(ctx, event); err != nil {
		log.Error("error while sending event after updating watering plan", "err", err, "watering_plan_id", prevWp.ID)
	}

	return nil
}

func (w *WateringPlanService) PreviewRoute(ctx context.Context, transporterID int32, trailerID *int32, clusterIDs []int32) (*entities.GeoJSON, error) {
	log := logger.GetLogger(ctx)
	transporter, err := w.vehicleRepo.GetByID(ctx, transporterID)
	if err != nil {
		log.Error("can't get selected transporter to preview route", "error", err, "vehicle_id", transporterID)
		return nil, service.NewError(service.NotFound, fmt.Sprintf("vehicle with id %d not found", transporterID))
	}

	var trailer *entities.Vehicle
	if trailerID != nil {
		trailer, err = w.vehicleRepo.GetByID(ctx, *trailerID)
		if err != nil {
			slog.Error("can't get selected trailer to preview route. route will be calculated without trailer", "error", err, "trailer_id", trailerID)
		}
	}

	clusters, err := w.clusterRepo.GetByIDs(ctx, clusterIDs)
	if err != nil {
		// when error, something is wrong with the db, else clusters should be an empty array
		log.Error("failed to get cluster by provided ids", "cluster_ids", clusterIDs)
		return nil, err
	}

	geoJSON, err := w.routingRepo.GenerateRoute(ctx, w.mergeVehicle(transporter, trailer), clusters)
	if err != nil {
		if errors.Is(err, storage.ErrUnknownVehicleType) {
			log.Debug("the vehicle type is not supported", "error", err, "vehicle_type", transporter.Type)
			return nil, service.NewError(service.NotFound, "vehicle type is not supported")
		}
		log.Debug("failed to generate route", "error", err)
		return nil, err
	}

	return geoJSON, nil
}

func (w *WateringPlanService) GetAll(ctx context.Context) ([]*entities.WateringPlan, error) {
	log := logger.GetLogger(ctx)
	plans, err := w.wateringPlanRepo.GetAll(ctx)
	if err != nil {
		log.Error("failed to fetch watering plans", "error", err)
		return nil, handleError(err)
	}

	return plans, nil
}

func (w *WateringPlanService) GetByID(ctx context.Context, id int32) (*entities.WateringPlan, error) {
	log := logger.GetLogger(ctx)
	got, err := w.wateringPlanRepo.GetByID(ctx, id)
	if err != nil {
		log.Error("failed to fetch watering plan by id", "error", err, "watering_plan_id", id)
		return nil, handleError(err)
	}

	return got, nil
}

func (w *WateringPlanService) Create(ctx context.Context, createWp *entities.WateringPlanCreate) (*entities.WateringPlan, error) {
	log := logger.GetLogger(ctx)
	if err := w.validator.Struct(createWp); err != nil {
		log.Debug("failed to validate struct from create watering plan", "error", err, "raw_watering_plan", fmt.Sprintf("%+v", createWp))
		return nil, service.NewError(service.BadRequest, errors.Join(err, errors.New("validation error")).Error())
	}

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
			log.Warn("failed to get trailer by id. will continue without trailer", "error", err)
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
		log.Error("failed to create watering plan", "error", err)
		return nil, handleError(err)
	}

	err = w.wateringPlanRepo.Update(ctx, created.ID, func(wp *entities.WateringPlan) (bool, error) {
		mergedVehicle := w.mergeVehicle(transporter, trailer)
		gpxURL, err := w.getGpxRouteURL(ctx, created.ID, mergedVehicle, treeClusters)
		if err == nil {
			wp.GpxURL = gpxURL
		}

		metadata, err := w.routingRepo.GenerateRouteInformation(ctx, mergedVehicle, treeClusters)
		if err == nil {
			wp.Distance = utils.P(metadata.Distance)
			wp.Duration = metadata.Time
			wp.RefillCount = metadata.Refills
		}

		return true, nil
	})

	if err != nil {
		log.Error("failed to apply generate gpx url to recently created watering plan", "error", err, "watering_plan_id", created.ID)
		return nil, handleError(err)
	}

	log.Info("watering plan created successfully", "watering_plan_id", created.ID)
	return created, nil
}

func (w *WateringPlanService) getGpxRouteURL(ctx context.Context, waterPlanID int32, vehicle *entities.Vehicle, clusters []*entities.TreeCluster) (string, error) {
	log := logger.GetLogger(ctx)
	r, err := w.routingRepo.GenerateRawGpxRoute(ctx, vehicle, clusters)
	if err != nil {
		log.Error("failed to generate gpx route", "error", err)
		return "", err
	}
	defer r.Close()

	objName := fmt.Sprintf("waterplan-%d.gpx", waterPlanID)

	var buf bytes.Buffer
	length, err := io.Copy(&buf, r)
	if err != nil {
		log.Error("error while reading gpx response", "error", err)
		return "", err
	}

	if err := w.gpxBucket.PutObject(ctx, objName, "application/gpx+xml;charset=UTF-8 ", length, &buf); err != nil {
		log.Error("failed to upload gpx object to bucket", "error", err, "bucket_name", viper.GetString("s3.route-gpx.bucket"), "obj_name", objName)
		return "", err
	}

	log.Info("gpx route successfully uploaded to s3 bucket", "obj_name", objName, "bucket_name", viper.GetString("s3.route-gpx.bucket"))
	return fmt.Sprintf("/v1/watering-plan/route/gpx/%s", objName), nil
}

func (w *WateringPlanService) GetGPXFileStream(ctx context.Context, objName string) (io.ReadSeekCloser, error) {
	log := logger.GetLogger(ctx)
	log.Debug("get gpx route object from bucket", "obj_name", objName, "bucket_name", viper.GetString("s3.route-gpx.bucket"))
	return w.gpxBucket.GetObject(ctx, objName)
}

func (w *WateringPlanService) Update(ctx context.Context, id int32, updateWp *entities.WateringPlanUpdate) (*entities.WateringPlan, error) {
	log := logger.GetLogger(ctx)
	if err := w.validator.Struct(updateWp); err != nil {
		log.Debug("failed to validate struct from update watering plan", "error", err, "raw_watering_plan", fmt.Sprintf("%+v", updateWp))
		return nil, service.NewError(service.BadRequest, errors.Join(err, errors.New("validation error")).Error())
	}

	prevWp, err := w.GetByID(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}

	// TODO: validate driver license

	if err := w.validateStatusDependentValues(ctx, updateWp); err != nil {
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

		mergedVehicle := w.mergeVehicle(transporter, trailer)
		if w.shouldUpdateGpx(prevWp, wp) {
			gpxURL, err := w.getGpxRouteURL(ctx, id, mergedVehicle, treeClusters)
			if err != nil {
				return false, handleError(err)
			}

			wp.GpxURL = gpxURL
		}

		metadata, err := w.routingRepo.GenerateRouteInformation(ctx, mergedVehicle, treeClusters)
		if err != nil {
			log.Error("failed to generate routing metadata", "error", err)
		} else {
			wp.Distance = utils.P(metadata.Distance)
			wp.Duration = metadata.Time
			wp.RefillCount = metadata.Refills
		}

		return true, nil
	})

	if err != nil {
		log.Error("failed to update watering plan")
		return nil, handleError(err)
	}

	log.Info("watering plan updated successfully", "watering_plan_id", id)
	if err := w.publishUpdateEvent(ctx, prevWp); err != nil {
		return nil, handleError(err)
	}

	return w.GetByID(ctx, id)
}

func (w *WateringPlanService) Delete(ctx context.Context, id int32) error {
	// _, err := w.wateringPlanRepo.GetByID(ctx, id)
	// if err != nil {
	// 	return handleError(err)
	// }

	log := logger.GetLogger(ctx)
	if err := w.wateringPlanRepo.Delete(ctx, id); err != nil {
		log.Error("failed to delete watering plan", "error", err, "watering_plan_id", id)
		return handleError(err)
	}

	log.Info("watering plan deleted successfully", "watering_plan_id", id)
	return nil
}

func (w *WateringPlanService) Ready() bool {
	return w.wateringPlanRepo != nil
}

func (w *WateringPlanService) shouldUpdateGpx(prevWp, newWp *entities.WateringPlan) bool {
	if len(prevWp.TreeClusters) != len(newWp.TreeClusters) {
		return true
	}

	if prevWp.Transporter.ID != newWp.Transporter.ID {
		return true
	}

	if (prevWp.Trailer == nil && newWp.Trailer != nil) || (prevWp.Trailer != nil && newWp.Trailer == nil) {
		return true
	}

	if prevWp.Trailer != nil && newWp.Trailer != nil && prevWp.Trailer.ID != newWp.Trailer.ID {
		return true
	}

	for i, prevWpTc := range prevWp.TreeClusters {
		if prevWpTc.ID != newWp.TreeClusters[i].ID {
			return true
		}
	}

	return false
}

func (w *WateringPlanService) fetchVehicle(ctx context.Context, vehicleID int32) (*entities.Vehicle, error) {
	log := logger.GetLogger(ctx)
	vehicle, err := w.vehicleRepo.GetByID(ctx, vehicleID)
	if err != nil {
		log.Error("failed to fetch vehicle by provided id", "error", err, "vehicle_id", vehicleID)
	}

	if vehicle == nil {
		return nil, service.NewError(service.NotFound, storage.ErrVehicleNotFound.Error())
	}

	return vehicle, nil
}

func (w *WateringPlanService) fetchTreeClusters(ctx context.Context, treeClusterIDs []*int32) ([]*entities.TreeCluster, error) {
	log := logger.GetLogger(ctx)
	clusters, err := w.getTreeClusters(ctx, treeClusterIDs)
	if err != nil {
		log.Error("failed to fetch tree cluster specified by requested ids", "cluster_ids", treeClusterIDs, "error", err)
		return nil, handleError(err)
	}
	if len(clusters) == 0 {
		log.Debug("requested tree cluster ids in watering plan are not found", "cluster_ids", treeClusterIDs, "error", err)
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
	log := logger.GetLogger(ctx)
	userIDStr := utils.Map(userIDs, func(id *uuid.UUID) string {
		return utils.UUIDToString(*id)
	})

	users, err := w.userRepo.GetByIDs(ctx, userIDStr)
	if err != nil {
		log.Error("failed to fetch users by id", "error", err, "user_ids", userIDStr)
		return handleError(err)
	}

	if len(users) == 0 {
		log.Debug("requested user ids in watering plan not found", "error", err, "user_ids", userIDStr)
		return storage.ErrUserNotFound
	}

	return nil
}

func (w *WateringPlanService) validateStatusDependentValues(ctx context.Context, entity *entities.WateringPlanUpdate) error {
	log := logger.GetLogger(ctx)

	// Set cancellation note to nothing if the current status is not fitting
	if entity.CancellationNote != "" && entity.Status != entities.WateringPlanStatusCanceled {
		log.Debug("cancellation note can only be set if watering plan is canceled")
		return service.NewError(service.BadRequest, "cancellation note can only be set if watering plan is canceled")
	}

	if entity.Status != entities.WateringPlanStatusFinished && len(entity.Evaluation) > 0 {
		log.Debug("evaluation values can only be set if the watering plan has been finished")
		return service.NewError(service.BadRequest, "evaluation values can only be set if the watering plan has been finished")
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

	var biggerWidth = transporter.Width
	if trailer.Width > transporter.Width {
		biggerWidth = trailer.Width
	}

	var biggerHeight = transporter.Height
	if trailer.Height > transporter.Height {
		biggerHeight = trailer.Height
	}

	return &entities.Vehicle{
		Width:         biggerWidth,
		Height:        biggerHeight,
		Length:        transporter.Length + trailer.Length,
		Weight:        transporter.Weight + trailer.Weight,
		WaterCapacity: transporter.WaterCapacity + trailer.WaterCapacity, // TODO: There may be a choice of transporter and trailer, but only the trailer will have water capacity should it be in use.
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

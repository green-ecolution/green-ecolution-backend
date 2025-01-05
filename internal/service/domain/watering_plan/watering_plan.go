package wateringplan

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/pkg/errors"

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
	validator        *validator.Validate
	eventManager    *worker.EventManager
}

func NewWateringPlanService(
	wateringPlanRepository storage.WateringPlanRepository,
	clusterRepository storage.TreeClusterRepository,
	vehicleRepository storage.VehicleRepository,
	userRepository storage.UserRepository,
	eventManager *worker.EventManager,
) service.WateringPlanService {
	return &WateringPlanService{
		wateringPlanRepo: wateringPlanRepository,
		clusterRepo:      clusterRepository,
		vehicleRepo:      vehicleRepository,
		userRepo:         userRepository,
		validator:        validator.New(),
		eventManager:    eventManager,
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
		return nil, service.NewError(service.BadRequest, errors.Wrap(err, "validation error").Error())
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
			return nil, err
		}
	}

	created, err := w.wateringPlanRepo.Create(ctx, func(wp *entities.WateringPlan) (bool, error) {
		wp.Date = createWp.Date
		wp.Description = createWp.Description
		wp.Transporter = transporter
		wp.Trailer = trailer
		wp.TreeClusters = treeClusters
		wp.UserIDs = createWp.UserIDs

		return true, nil
	})

	if err != nil {
		return nil, handleError(err)
	}

	return created, nil
}

func (w *WateringPlanService) Update(ctx context.Context, id int32, updateWp *entities.WateringPlanUpdate) (*entities.WateringPlan, error) {
	if err := w.validator.Struct(updateWp); err != nil {
		return nil, service.NewError(service.BadRequest, errors.Wrap(err, "validation error").Error())
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

		return true, nil
	})

	if err != nil {
		return nil, handleError(err)
	}

	w.publishUpdateEvent(ctx, prevWp)
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

func handleError(err error) error {
	if errors.Is(err, storage.ErrEntityNotFound) {
		return service.NewError(service.NotFound, storage.ErrWateringPlanNotFound.Error())
	}

	return service.NewError(service.InternalError, err.Error())
}

package wateringplan

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
)

type WateringPlanService struct {
	wateringPlanRepo storage.WateringPlanRepository
	clusterRepo      storage.TreeClusterRepository
	vehicleRepo      storage.VehicleRepository
	validator        *validator.Validate
	statusSchedular  *StatusSchedular
}

func NewWateringPlanService(
	wateringPlanRepository storage.WateringPlanRepository,
	clusterRepository storage.TreeClusterRepository,
	vehicleRepository storage.VehicleRepository,
) service.WateringPlanService {
	return &WateringPlanService{
		wateringPlanRepo: wateringPlanRepository,
		clusterRepo:      clusterRepository,
		vehicleRepo:      vehicleRepository,
		validator:        validator.New(),
		statusSchedular:  &StatusSchedular{wateringPlanRepo: wateringPlanRepository},
	}
}

func (w *WateringPlanService) RunStatusSchedular(ctx context.Context, interval time.Duration) {
	w.statusSchedular.RunStatusSchedular(ctx, interval)
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

	// TODO: get users
	// TODO: calculare distance

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

	// TODO: get users
	// TODO: calculate distance

	if err := w.validateStatusDependentValues(updateWp); err != nil {
		return nil, err
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

		return true, nil
	})

	if err != nil {
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

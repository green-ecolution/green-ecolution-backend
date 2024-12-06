package wateringplan

import (
	"context"

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
	}
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

	// TODO: get clusters, vehicles, users
	// TODO: calculate required water
	// TODO: calculare distance

	transporter, err := w.getVehicle(ctx, createWp.TransporterID)
	if err != nil {
		return nil, service.NewError(service.NotFound, storage.ErrVehicleNotFound.Error())
	}

	trailer, err := w.getVehicle(ctx, createWp.TrailerID)
	if err != nil && !errors.Is(err, storage.ErrVehicleNotFound){
		return nil, service.NewError(service.InternalError, err.Error())
	}

	created, err := w.wateringPlanRepo.Create(ctx, func(wp *entities.WateringPlan) (bool, error) {
		wp.Date = createWp.Date
		wp.Description = createWp.Description
		wp.Transporter = transporter
		wp.Trailer = trailer

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

	// TODO: get clusters & users
	// TODO: calculate required water
	// TODO: calculare distance

	// transporter, err := w.getVehicle(ctx, updateWp.TransporterID)
	// if err != nil {
	// 	return nil, service.NewError(service.NotFound, storage.ErrVehicleNotFound.Error())
	// }

	// trailer, err := w.getVehicle(ctx, updateWp.TrailerID)
	// if err != nil && !errors.Is(err, storage.ErrVehicleNotFound){
	// 	return nil, service.NewError(service.InternalError, err.Error())
	// }

	err := w.wateringPlanRepo.Update(ctx, id, func(wp *entities.WateringPlan) (bool, error) {
		wp.Date = updateWp.Date
		wp.Description = updateWp.Description
		// wp.Transporter = transporter
		// wp.Trailer = trailer

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

func handleError(err error) error {
	if errors.Is(err, storage.ErrEntityNotFound) {
		return service.NewError(service.NotFound, storage.ErrWateringPlanNotFound.Error())
	}

	return service.NewError(service.InternalError, err.Error())
}

func (w *WateringPlanService) getVehicle(ctx context.Context, id *int32) (*entities.Vehicle, error) {
	var err error
	vehicle, err := w.vehicleRepo.GetByID(ctx, *id)
	if err != nil {
		return nil, err
	}

	return vehicle, nil
}

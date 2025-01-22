package vehicle

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
)

type VehicleService struct {
	vehicleRepo storage.VehicleRepository
	validator   *validator.Validate
}

func NewVehicleService(vehicleRepository storage.VehicleRepository) service.VehicleService {
	return &VehicleService{
		vehicleRepo: vehicleRepository,
		validator:   validator.New(),
	}
}

func (v *VehicleService) GetAll(ctx context.Context) ([]*entities.Vehicle, error) {
	log := logger.GetLogger(ctx)
	vehicles, err := v.vehicleRepo.GetAll(ctx)
	if err != nil {
		log.Error("failed to fetch vehicles", "error", err)
		return nil, handleError(err)
	}

	return vehicles, nil
}

func (v *VehicleService) GetAllByType(ctx context.Context, vehicleType entities.VehicleType) ([]*entities.Vehicle, error) {
	log := logger.GetLogger(ctx)
	vehicles, err := v.vehicleRepo.GetAllByType(ctx, vehicleType)
	if err != nil {
		log.Error("failed to fetch vehicles by a type", "error", err, "vehicle_type", vehicleType)
		return nil, handleError(err)
	}

	return vehicles, nil
}

func (v *VehicleService) GetByID(ctx context.Context, id int32) (*entities.Vehicle, error) {
	log := logger.GetLogger(ctx)
	got, err := v.vehicleRepo.GetByID(ctx, id)
	if err != nil {
		log.Error("failed to fetch vehicle by id", "error", err, "vehicle_id", id)
		return nil, handleError(err)
	}

	if got == nil {
		return nil, service.NewError(service.NotFound, storage.ErrVehicleNotFound.Error()) // TODO: change to service error
	}

	return got, nil
}

func (v *VehicleService) GetByPlate(ctx context.Context, plate string) (*entities.Vehicle, error) {
	log := logger.GetLogger(ctx)
	got, err := v.vehicleRepo.GetByPlate(ctx, plate)
	if err != nil {
		log.Error("failed to fetch vehicle by plate", "error", err, "vehicle_plate", plate)
		return nil, handleError(err)
	}

	return got, nil
}

func (v *VehicleService) Create(ctx context.Context, createData *entities.VehicleCreate) (*entities.Vehicle, error) {
	log := logger.GetLogger(ctx)
	if err := v.validator.Struct(createData); err != nil {
		log.Debug("failed to validate struct from create vehicle", "error", err, "raw_vehicle", fmt.Sprintf("%+v", createData))
		return nil, service.NewError(service.BadRequest, errors.Wrap(err, "validation error").Error())
	}

	if isTaken, err := v.isVehicleNumberPlateTaken(ctx, createData.NumberPlate); err != nil {
		log.Error("failed to request if vehicle plate is already taken", "error", err, "vehicle_plate", createData.NumberPlate)
		return nil, err
	} else if isTaken {
		log.Debug("requested number plate is already taken", "vehicle_plate", createData.NumberPlate)
		return nil, service.NewError(service.BadRequest, errors.New("number plate is already taken").Error())
	}

	created, err := v.vehicleRepo.Create(ctx, func(vh *entities.Vehicle) (bool, error) {
		vh.NumberPlate = createData.NumberPlate
		vh.Description = createData.Description
		vh.WaterCapacity = createData.WaterCapacity
		vh.Status = createData.Status
		vh.Type = createData.Type
		vh.Height = createData.Height
		vh.Length = createData.Length
		vh.Width = createData.Width
		vh.Model = createData.Model
		vh.Weight = createData.Weight
		vh.DrivingLicense = createData.DrivingLicense

		return true, nil
	})
	if err != nil {
		log.Error("failed to create vehicle", "error", err)
		return nil, handleError(err)
	}

	log.Info("vehicle created successfully", "vehicle_id", created.ID)
	return created, nil
}

func (v *VehicleService) Update(ctx context.Context, id int32, updateData *entities.VehicleUpdate) (*entities.Vehicle, error) {
	log := logger.GetLogger(ctx)
	if err := v.validator.Struct(updateData); err != nil {
		log.Debug("failed to validate struct from update vehicle", "error", err, "raw_vehicle", fmt.Sprintf("%+v", updateData))
		return nil, service.NewError(service.BadRequest, errors.Wrap(err, "validation error").Error())
	}

	oldValue, err := v.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if oldValue.NumberPlate != updateData.NumberPlate {
		if isTaken, err := v.isVehicleNumberPlateTaken(ctx, updateData.NumberPlate); err != nil {
			log.Error("failed to request if vehicle plate is already taken", "error", err, "vehicle_plate", updateData.NumberPlate)
			return nil, err
		} else if isTaken {
			log.Debug("requested number plate is already taken", "vehicle_plate", updateData.NumberPlate)
			return nil, service.NewError(service.BadRequest, errors.New("number plate is already taken").Error()) // TODO: move to svc error
		}
	}

	err = v.vehicleRepo.Update(ctx, id, func(vh *entities.Vehicle) (bool, error) {
		vh.NumberPlate = updateData.NumberPlate
		vh.Description = updateData.Description
		vh.WaterCapacity = updateData.WaterCapacity
		vh.Status = updateData.Status
		vh.Type = updateData.Type
		vh.Height = updateData.Height
		vh.Length = updateData.Length
		vh.Width = updateData.Width
		vh.Model = updateData.Model
		vh.Weight = updateData.Weight
		vh.DrivingLicense = updateData.DrivingLicense

		return true, nil
	})

	if err != nil {
		log.Error("failed to update vehicle", "error", err, "vehicle_id", id)
		return nil, handleError(err)
	}

	log.Info("vehicle updated successfully", "vehicle_id", id)
	return v.GetByID(ctx, id)
}

func (v *VehicleService) Delete(ctx context.Context, id int32) error {
	log := logger.GetLogger(ctx)
	if _, err := v.vehicleRepo.GetByID(ctx, id); err != nil {
		return handleError(err)
	}

	if err := v.vehicleRepo.Delete(ctx, id); err != nil {
		log.Error("failed to delete vehicle", "error", err, "vehicle_id", id)
		return handleError(err)
	}

	log.Info("vehicle deleted successfully", "vehicle_id", id)
	return nil
}

func (v *VehicleService) Ready() bool {
	return v.vehicleRepo != nil
}

func handleError(err error) error {
	if errors.Is(err, storage.ErrEntityNotFound) {
		return service.NewError(service.NotFound, storage.ErrVehicleNotFound.Error())
	}

	return service.NewError(service.InternalError, err.Error())
}

func (v *VehicleService) isVehicleNumberPlateTaken(ctx context.Context, plate string) (bool, error) {
	existingVehicle, err := v.vehicleRepo.GetByPlate(ctx, plate)
	if err != nil && !errors.Is(err, storage.ErrEntityNotFound) {
		return false, handleError(err)
	}
	return existingVehicle != nil, nil
}

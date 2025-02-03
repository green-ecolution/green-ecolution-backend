package vehicle

import (
	"context"
	"errors"
	"fmt"

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

func (v *VehicleService) GetAll(ctx context.Context, provider string) ([]*entities.Vehicle, error) {
	log := logger.GetLogger(ctx)
	var vehicles []*entities.Vehicle
	var err error

	if provider != "" {
		vehicles, err = v.vehicleRepo.GetAllByProvider(ctx, provider)
	} else {
		vehicles, err = v.vehicleRepo.GetAll(ctx)
	}
	if err != nil {
		log.Debug("failed to fetch vehicles", "error", err)
		return nil, service.MapError(ctx, err, service.ErrorLogEntityNotFound)
	}

	return vehicles, nil
}

func (v *VehicleService) GetAllByType(ctx context.Context, vehicleType entities.VehicleType) ([]*entities.Vehicle, error) {
	log := logger.GetLogger(ctx)
	vehicles, err := v.vehicleRepo.GetAllByType(ctx, vehicleType)
	if err != nil {
		log.Debug("failed to fetch vehicles by a type", "error", err, "vehicle_type", vehicleType)
		return nil, service.MapError(ctx, err, service.ErrorLogEntityNotFound)
	}

	return vehicles, nil
}

func (v *VehicleService) GetByID(ctx context.Context, id int32) (*entities.Vehicle, error) {
	log := logger.GetLogger(ctx)
	got, err := v.vehicleRepo.GetByID(ctx, id)
	if err != nil {
		log.Debug("failed to fetch vehicle by id", "error", err, "vehicle_id", id)
		return nil, service.MapError(ctx, err, service.ErrorLogEntityNotFound)
	}

	return got, nil
}

func (v *VehicleService) GetByPlate(ctx context.Context, plate string) (*entities.Vehicle, error) {
	log := logger.GetLogger(ctx)
	got, err := v.vehicleRepo.GetByPlate(ctx, plate)
	if err != nil {
		log.Debug("failed to fetch vehicle by plate", "error", err, "vehicle_plate", plate)
		return nil, service.MapError(ctx, err, service.ErrorLogEntityNotFound)
	}

	return got, nil
}

func (v *VehicleService) Create(ctx context.Context, createData *entities.VehicleCreate) (*entities.Vehicle, error) {
	log := logger.GetLogger(ctx)
	if err := v.validator.Struct(createData); err != nil {
		log.Debug("failed to validate struct from create vehicle", "error", err, "raw_vehicle", fmt.Sprintf("%+v", createData))
		return nil, service.MapError(ctx, errors.Join(err, service.ErrValidation), service.ErrorLogValidation)
	}

	if isTaken, err := v.isVehicleNumberPlateTaken(ctx, createData.NumberPlate); err != nil {
		log.Debug("failed to request if vehicle plate is already taken", "error", err, "vehicle_plate", createData.NumberPlate)
		return nil, service.MapError(ctx, err, service.ErrorLogAll)
	} else if isTaken {
		log.Debug("requested number plate is already taken", "vehicle_plate", createData.NumberPlate)
		return nil, service.ErrVehiclePlateTaken
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
		vh.Provider = createData.Provider
		vh.AdditionalInfo = createData.AdditionalInfo

		return true, nil
	})
	if err != nil {
		log.Error("failed to create vehicle", "error", err)
		return nil, service.MapError(ctx, err, service.ErrorLogAll)
	}

	log.Info("vehicle created successfully", "vehicle_id", created.ID)
	return created, nil
}

func (v *VehicleService) Update(ctx context.Context, id int32, updateData *entities.VehicleUpdate) (*entities.Vehicle, error) {
	log := logger.GetLogger(ctx)
	if err := v.validator.Struct(updateData); err != nil {
		log.Debug("failed to validate struct from update vehicle", "error", err, "raw_vehicle", fmt.Sprintf("%+v", updateData))
		return nil, service.MapError(ctx, errors.Join(err, service.ErrValidation), service.ErrorLogValidation)
	}

	oldValue, err := v.GetByID(ctx, id)
	if err != nil {
		log.Debug("failed to get already existing vehicle from store", "error", err, "vehicle_id", id)
		return nil, service.MapError(ctx, err, service.ErrorLogEntityNotFound)
	}

	if oldValue.NumberPlate != updateData.NumberPlate {
		if isTaken, err := v.isVehicleNumberPlateTaken(ctx, updateData.NumberPlate); err != nil {
			log.Debug("failed to request if vehicle plate is already taken", "error", err, "vehicle_plate", updateData.NumberPlate)
			return nil, service.MapError(ctx, err, service.ErrorLogAll)
		} else if isTaken {
			log.Debug("requested number plate is already taken", "vehicle_plate", updateData.NumberPlate)
			return nil, service.ErrVehiclePlateTaken
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
		vh.Provider = updateData.Provider
		vh.AdditionalInfo = updateData.AdditionalInfo

		return true, nil
	})

	if err != nil {
		log.Debug("failed to update vehicle", "error", err, "vehicle_id", id)
		return nil, service.MapError(ctx, err, service.ErrorLogAll)
	}

	log.Info("vehicle updated successfully", "vehicle_id", id)
	return v.GetByID(ctx, id)
}

func (v *VehicleService) Delete(ctx context.Context, id int32) error {
	log := logger.GetLogger(ctx)
	if _, err := v.vehicleRepo.GetByID(ctx, id); err != nil {
		log.Debug("failed to get vehicle by id in delete request", "error", err, "vehicle_id", id)
		return service.MapError(ctx, err, service.ErrorLogEntityNotFound)
	}

	if err := v.vehicleRepo.Delete(ctx, id); err != nil {
		log.Debug("failed to delete vehicle", "error", err, "vehicle_id", id)
		return service.MapError(ctx, err, service.ErrorLogAll)
	}

	log.Info("vehicle deleted successfully", "vehicle_id", id)
	return nil
}

func (v *VehicleService) Ready() bool {
	return v.vehicleRepo != nil
}

func (v *VehicleService) isVehicleNumberPlateTaken(ctx context.Context, plate string) (bool, error) {
	existingVehicle, err := v.vehicleRepo.GetByPlate(ctx, plate)
	var entityNotFoundErr storage.ErrEntityNotFound
	if err != nil && !errors.As(err, &entityNotFoundErr) {
		return false, err
	}
	return existingVehicle != nil, nil
}

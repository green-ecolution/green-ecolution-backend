package vehicle

import (
	"context"

	"github.com/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
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
	vehicles, err := v.vehicleRepo.GetAll(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	return vehicles, nil
}

func (v *VehicleService) GetByID(ctx context.Context, id int32) (*entities.Vehicle, error) {
	got, err := v.vehicleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}

	return got, nil
}

func (v *VehicleService) GetByPlate(ctx context.Context, plate string) (*entities.Vehicle, error) {
	got, err := v.vehicleRepo.GetByPlate(ctx, plate)
	if err != nil {
		return nil, handleError(err)
	}

	return got, nil
}

func (v *VehicleService) Create(ctx context.Context, createData *entities.VehicleCreate) (*entities.Vehicle, error) {
	if err := v.validator.Struct(createData); err != nil {
		return nil, service.NewError(service.BadRequest, errors.Wrap(err, "validation error").Error())
	}

	if isTaken, err := v.isVehicleNumberPlateTaken(ctx, createData.NumberPlate); err != nil {
		return nil, err
	} else if isTaken {
		return nil, service.NewError(service.BadRequest, errors.New("Number plate is already taken").Error())
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
		vh.DrivingLicense = createData.DrivingLicense

		return true, nil
	})

	if err != nil {
		return nil, handleError(err)
	}

	return created, nil
}

func (v *VehicleService) Update(ctx context.Context, id int32, updateData *entities.VehicleUpdate) (*entities.Vehicle, error) {
	if err := v.validator.Struct(updateData); err != nil {
		return nil, service.NewError(service.BadRequest, errors.Wrap(err, "validation error").Error())
	}

	oldValue, err := v.vehicleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}

	if oldValue.NumberPlate != updateData.NumberPlate {
		if isTaken, err := v.isVehicleNumberPlateTaken(ctx, updateData.NumberPlate); err != nil {
			return nil, err
		} else if isTaken {
			return nil, service.NewError(service.BadRequest, errors.New("Number plate is already taken").Error())
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
		vh.DrivingLicense = updateData.DrivingLicense

		return true, nil
	})

	if err != nil {
		return nil, handleError(err)
	}

	return v.GetByID(ctx, id)
}

func (v *VehicleService) Delete(ctx context.Context, id int32) error {
	_, err := v.vehicleRepo.GetByID(ctx, id)
	if err != nil {
		return handleError(err)
	}

	err = v.vehicleRepo.Delete(ctx, id)
	if err != nil {
		return handleError(err)
	}

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

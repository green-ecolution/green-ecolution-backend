package vehicle

import (
	"context"
	"log/slog"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper"

	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	store "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
)

type VehicleRepository struct {
	store *store.Store
	VehicleRepositoryMappers
}

type VehicleRepositoryMappers struct {
	mapper mapper.InternalVehicleRepoMapper
}

func NewVehicleRepositoryMappers(vMapper mapper.InternalVehicleRepoMapper) VehicleRepositoryMappers {
	return VehicleRepositoryMappers{
		mapper: vMapper,
	}
}

func NewVehicleRepository(s *store.Store, mappers VehicleRepositoryMappers) storage.VehicleRepository {
	return &VehicleRepository{
		store:                    s,
		VehicleRepositoryMappers: mappers,
	}
}

func WithNumberPlate(numberPlate string) entities.EntityFunc[entities.Vehicle] {
	return func(v *entities.Vehicle) {
		slog.Debug("updating number", "number", numberPlate)
		v.NumberPlate = numberPlate
	}
}

func WithDescription(description string) entities.EntityFunc[entities.Vehicle] {
	return func(v *entities.Vehicle) {
		slog.Debug("updating description", "description", description)
		v.Description = description
	}
}

func WithWaterCapacity(waterCapacity float64) entities.EntityFunc[entities.Vehicle] {
	return func(v *entities.Vehicle) {
		slog.Debug("updating water capacity", "water capacity", waterCapacity)
		v.WaterCapacity = waterCapacity
	}
}

func WithVehicleType(vehicleType entities.VehicleType) entities.EntityFunc[entities.Vehicle] {
	return func(v *entities.Vehicle) {
		slog.Debug("updating vehicle type", "vehicle type", vehicleType)
		v.Type = vehicleType
	}
}

func WithVehicleStatus(vehicleStatus entities.VehicleStatus) entities.EntityFunc[entities.Vehicle] {
	return func(v *entities.Vehicle) {
		slog.Debug("updating vehicle status", "vehicle status", vehicleStatus)
		v.Status = vehicleStatus
	}
}

func WithModel(model string) entities.EntityFunc[entities.Vehicle] {
	return func(v *entities.Vehicle) {
		slog.Debug("updating model", "model", model)
		v.Model = model
	}
}

func WithDrivingLicense(drivingLicense entities.DrivingLicense) entities.EntityFunc[entities.Vehicle] {
	return func(v *entities.Vehicle) {
		slog.Debug("updating driving license", "driving license", drivingLicense)
		v.DrivingLicense = drivingLicense
	}
}

func WithHeight(height float64) entities.EntityFunc[entities.Vehicle] {
	return func(v *entities.Vehicle) {
		slog.Debug("updating height", "height", height)
		v.Height = height
	}
}

func WithWidth(width float64) entities.EntityFunc[entities.Vehicle] {
	return func(v *entities.Vehicle) {
		slog.Debug("updating width", "width", width)
		v.Width = width
	}
}

func WithLength(length float64) entities.EntityFunc[entities.Vehicle] {
	return func(v *entities.Vehicle) {
		slog.Debug("updating length", "length", length)
		v.Length = length
	}
}

func (r *VehicleRepository) Delete(ctx context.Context, id int32) error {
	_, err := r.store.DeleteVehicle(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

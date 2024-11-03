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
	s.SetEntityType(store.Vehicle)
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

func (r *VehicleRepository) Delete(ctx context.Context, id int32) error {
	rowID, err := r.store.DeleteVehicle(ctx, id)
	if err != nil {
		return err
	}

	if rowID != id || rowID == 0 {
		return storage.ErrVehicleNotFound
	}

	return nil
}

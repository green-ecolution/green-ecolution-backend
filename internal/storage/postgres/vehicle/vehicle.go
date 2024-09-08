package vehicle

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	. "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
)

type VehicleRepository struct {
	store *Store
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

func NewVehicleRepository(store *Store, mappers VehicleRepositoryMappers) storage.VehicleRepository {
	return &VehicleRepository{
		store:                    store,
		VehicleRepositoryMappers: mappers,
	}
}

func (r *VehicleRepository) GetAll(ctx context.Context) ([]*entities.Vehicle, error) {
	rows, err := r.store.GetAllVehicles(ctx)
	if err != nil {
		return nil, err
	}

	return r.mapper.FromSqlList(rows), nil
}

func (r *VehicleRepository) GetByID(ctx context.Context, id int32) (*entities.Vehicle, error) {
	row, err := r.store.GetVehicleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return r.mapper.FromSql(row), nil
}

func (r *VehicleRepository) GetByPlate(ctx context.Context, plate string) (*entities.Vehicle, error) {
	row, err := r.store.GetVehicleByPlate(ctx, plate)
	if err != nil {
		return nil, err
	}

	return r.mapper.FromSql(row), nil
}

func (r *VehicleRepository) Create(ctx context.Context, vehicle *entities.CreateVehicle) (*entities.Vehicle, error) {
	params := &sqlc.CreateVehicleParams{
		NumberPlate:   vehicle.NumberPlate,
		Description:   vehicle.Description,
		WaterCapacity: vehicle.WaterCapacity,
	}

	id, err := r.store.CreateVehicle(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *VehicleRepository) Update(ctx context.Context, vehicle *entities.UpdateVehicle) (*entities.Vehicle, error) {
	prev, err := r.GetByID(ctx, vehicle.ID)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	params := &sqlc.UpdateVehicleParams{
		ID:            vehicle.ID,
		NumberPlate:   utils.CompareAndUpdate(prev.NumberPlate, vehicle.NumberPlate),
		Description:   utils.CompareAndUpdate(prev.Description, vehicle.Description),
		WaterCapacity: utils.CompareAndUpdate(prev.WaterCapacity, vehicle.WaterCapacity),
	}

	if err = r.store.UpdateVehicle(ctx, params); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, vehicle.ID)
}

func (r *VehicleRepository) Delete(ctx context.Context, id int32) error {
	return r.store.DeleteVehicle(ctx, id)
}

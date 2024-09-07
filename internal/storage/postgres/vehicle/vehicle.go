package vehicle

import (
	"context"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

type VehicleRepository struct {
	querier sqlc.Querier
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

func NewVehicleRepository(querier sqlc.Querier, mappers VehicleRepositoryMappers) storage.VehicleRepository {
	return &VehicleRepository{
		querier:                  querier,
		VehicleRepositoryMappers: mappers,
	}
}

func (r *VehicleRepository) GetAll(ctx context.Context) ([]*entities.Vehicle, error) {
	rows, err := r.querier.GetAllVehicles(ctx)
	if err != nil {
		return nil, err
	}

	return r.mapper.FromSqlList(rows), nil
}

func (r *VehicleRepository) GetByID(ctx context.Context, id int32) (*entities.Vehicle, error) {
	row, err := r.querier.GetVehicleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return r.mapper.FromSql(row), nil
}

func (r *VehicleRepository) GetByPlate(ctx context.Context, plate string) (*entities.Vehicle, error) {
	row, err := r.querier.GetVehicleByPlate(ctx, plate)
	if err != nil {
		return nil, err
	}

	return r.mapper.FromSql(row), nil
}

func (r *VehicleRepository) Create(ctx context.Context, vehicle *entities.Vehicle) (*entities.Vehicle, error) {
	params := &sqlc.CreateVehicleParams{
		NumberPlate:   vehicle.NumberPlate,
		Description:   vehicle.Description,
		WaterCapacity: vehicle.WaterCapacity,
	}

	id, err := r.querier.CreateVehicle(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *VehicleRepository) Update(ctx context.Context, vehicle *entities.Vehicle) (*entities.Vehicle, error) {
	params := &sqlc.UpdateVehicleParams{
		ID:            vehicle.ID,
		NumberPlate:   vehicle.NumberPlate,
		Description:   vehicle.Description,
		WaterCapacity: vehicle.WaterCapacity,
	}

	err := r.querier.UpdateVehicle(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, vehicle.ID)
}

func (r *VehicleRepository) Delete(ctx context.Context, id int32) error {
	return r.querier.DeleteVehicle(ctx, id)
}

package vehicle

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

func (r *VehicleRepository) GetAll(ctx context.Context) ([]*entities.Vehicle, error) {
	log := logger.GetLogger(ctx)
	rows, err := r.store.GetAllVehicles(ctx)
	if err != nil {
		log.Debug("failed to get vehicle entities in db", "error", err)
		return nil, r.store.MapError(err, sqlc.Vehicle{})
	}

	return r.mapFromList(ctx, rows)
}

func (r *VehicleRepository) GetAllByProvider(ctx context.Context, provider string) ([]*entities.Vehicle, error) {
	log := logger.GetLogger(ctx)
	rows, err := r.store.GetAllVehiclesByProvider(ctx, &provider)
	if err != nil {
		log.Debug("failed to get vehicle entities in db", "error", err)
		return nil, r.store.MapError(err, sqlc.Vehicle{})
	}

	return r.mapFromList(ctx, rows)
}

func (r *VehicleRepository) GetAllByType(ctx context.Context, vehicleType entities.VehicleType) ([]*entities.Vehicle, error) {
	log := logger.GetLogger(ctx)
	rows, err := r.store.GetAllVehiclesByType(ctx, sqlc.VehicleType(vehicleType))
	if err != nil {
		log.Debug("failed to get vehicle entities by provides type in db", "error", err, "vehicle_type", vehicleType)
		return nil, r.store.MapError(err, sqlc.Vehicle{})
	}

	return r.mapFromList(ctx, rows)
}

func (r *VehicleRepository) GetByID(ctx context.Context, id int32) (*entities.Vehicle, error) {
	log := logger.GetLogger(ctx)
	row, err := r.store.GetVehicleByID(ctx, id)
	if err != nil {
		log.Debug("failed to get vehicle entity by provided id", "error", err, "vehicle_id", id)
		return nil, r.store.MapError(err, sqlc.Vehicle{})
	}

	return r.mapFromRow(ctx, row)
}

func (r *VehicleRepository) GetByPlate(ctx context.Context, plate string) (*entities.Vehicle, error) {
	log := logger.GetLogger(ctx)
	row, err := r.store.GetVehicleByPlate(ctx, plate)
	if err != nil {
		log.Debug("failed to get vehicle entity by given plate", "error", err, "vehicle_plate", plate)
		return nil, r.store.MapError(err, sqlc.Vehicle{})
	}

	return r.mapFromRow(ctx, row)
}

func (r *VehicleRepository) mapFromRow(ctx context.Context, rows *sqlc.Vehicle) (*entities.Vehicle, error) {
	log := logger.GetLogger(ctx)
	vehicles, err := r.mapper.FromSql(rows)
	if err != nil {
		log.Debug("failed to convert entity", "error", err)
		return nil, err
	}

	return vehicles, nil
}

func (r *VehicleRepository) mapFromList(ctx context.Context, rows []*sqlc.Vehicle) ([]*entities.Vehicle, error) {
	log := logger.GetLogger(ctx)
	vehicles, err := r.mapper.FromSqlList(rows)
	if err != nil {
		log.Debug("failed to convert entity", "error", err)
		return nil, err
	}

	return vehicles, nil
}

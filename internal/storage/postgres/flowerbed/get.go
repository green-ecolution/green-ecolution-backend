package flowerbed

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func (r *FlowerbedRepository) GetAll(ctx context.Context) ([]*entities.Flowerbed, error) {
	row, err := r.store.GetAllFlowerbeds(ctx)
	if err != nil {
		return nil, err
	}

	data := r.mapper.FromSqlList(row)
	for _, f := range data {
		f.Sensor, err = r.GetSensorByFlowerbedID(ctx, f.ID)
		if err != nil && !errors.Is(err, storage.ErrSensorNotFound) {
			return nil, err
		}

		f.Images, err = r.GetAllImagesByID(ctx, f.ID)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}

		f.Region, err = r.GetRegionByFlowerbedID(ctx, f.ID)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (r *FlowerbedRepository) GetByID(ctx context.Context, id int32) (*entities.Flowerbed, error) {
	row, err := r.store.GetFlowerbedByID(ctx, id)
	if err != nil {
		return nil, err
	}

	data := r.mapper.FromSql(row)

	data.Sensor, err = r.GetSensorByFlowerbedID(ctx, id)
	if err != nil && !errors.Is(err, storage.ErrSensorNotFound) {
		return nil, err
	}

	data.Images, err = r.GetAllImagesByID(ctx, id)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	data.Region, err = r.GetRegionByFlowerbedID(ctx, id)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *FlowerbedRepository) GetAllImagesByID(ctx context.Context, flowerbedID int32) ([]*entities.Image, error) {
	if err := r.flowerbedIDExists(ctx, flowerbedID); err != nil {
		return nil, err
	}

	row, err := r.store.GetAllImagesByFlowerbedID(ctx, flowerbedID)
	if err != nil {
		return nil, err
	}

	return r.imgMapper.FromSqlList(row), nil
}

func (r *FlowerbedRepository) GetSensorByFlowerbedID(ctx context.Context, flowerbedID int32) (*entities.Sensor, error) {
	if err := r.flowerbedIDExists(ctx, flowerbedID); err != nil {
		return nil, err
	}

	row, err := r.store.GetSensorByFlowerbedID(ctx, flowerbedID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrSensorNotFound
		}
		return nil, err
	}

	return r.sensorMapper.FromSql(row), nil
}

func (r *FlowerbedRepository) GetRegionByFlowerbedID(ctx context.Context, flowerbedID int32) (*entities.Region, error) {
	if err := r.flowerbedIDExists(ctx, flowerbedID); err != nil {
		return nil, err
	}

	row, err := r.store.GetRegionByFlowerbedID(ctx, flowerbedID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrRegionNotFound
		}
		return nil, err
	}

	return r.regionMapper.FromSql(row), nil
}

func (r *FlowerbedRepository) flowerbedIDExists(ctx context.Context, id int32) error {
	_, err := r.store.GetFlowerbedByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.ErrFlowerbedNotFound
		}
		return err
	}

	return nil
}

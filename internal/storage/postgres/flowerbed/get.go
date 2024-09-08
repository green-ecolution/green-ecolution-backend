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
		return nil, r.store.HandleError(err)
	}

	mapped := r.mapper.FromSqlList(row)
	for _, f := range mapped {
		if err := r.mapSensorAndImages(ctx, f); err != nil {
			return nil, r.store.HandleError(err)
		}
	}

	return mapped, nil
}

func (r *FlowerbedRepository) GetByID(ctx context.Context, id int32) (*entities.Flowerbed, error) {
	row, err := r.store.GetFlowerbedByID(ctx, id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	f := r.mapper.FromSql(row)
	if err := r.mapSensorAndImages(ctx, f); err != nil {
		return nil, r.store.HandleError(err)
	}

	return f, nil
}

func (r *FlowerbedRepository) GetAllImagesByID(ctx context.Context, flowerbedID int32) ([]*entities.Image, error) {
	row, err := r.store.GetAllImagesByFlowerbedID(ctx, flowerbedID)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	return r.imgMapper.FromSqlList(row), nil
}

func (r *FlowerbedRepository) GetSensorByFlowerbedID(ctx context.Context, flowerbedID int32) (*entities.Sensor, error) {
	row, err := r.store.GetSensorByFlowerbedID(ctx, flowerbedID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrSensorNotFound
		} else {
			return nil, r.store.HandleError(err)
		}
	}

	return r.sensorMapper.FromSql(row), nil
}

// Map sensor and images entity to domain flowerbed
func (r *FlowerbedRepository) mapSensorAndImages(ctx context.Context, f *entities.Flowerbed) error {
	if err := r.mapImages(ctx, f); err != nil {
		return err
	}

	if err := r.mapSensor(ctx, f); err != nil {
		return err
	}

	return nil
}

func (r *FlowerbedRepository) mapImages(ctx context.Context, f *entities.Flowerbed) error {
	images, err := r.GetAllImagesByID(ctx, f.ID)
	if err != nil {
		return r.store.HandleError(err)
	}
	f.Images = images
	return nil
}

func (r *FlowerbedRepository) mapSensor(ctx context.Context, f *entities.Flowerbed) error {
	sensor, err := r.GetSensorByFlowerbedID(ctx, f.ID)
	if err != nil {
		return r.store.HandleError(err)
	}
	f.Sensor = sensor
	return nil
}

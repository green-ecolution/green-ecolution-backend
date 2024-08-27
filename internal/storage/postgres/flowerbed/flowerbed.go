package flowerbed

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/flowerbed/mapper"
	imgMapper "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/image/mapper"
	sensorMapper "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/sensor/mapper"
	"github.com/pkg/errors"
)

type FlowerbedRepository struct {
	querier sqlc.Querier
	FlowerbedMappers
}

type FlowerbedMappers struct {
	mapper       mapper.InternalFlowerbedRepoMapper
	imgMapper    imgMapper.InternalImageRepoMapper
	sensorMapper sensorMapper.InternalSensorRepoMapper
}

func NewFlowerbedMappers(fMapper mapper.InternalFlowerbedRepoMapper, iMapper imgMapper.InternalImageRepoMapper, sMapper sensorMapper.InternalSensorRepoMapper) FlowerbedMappers {
	return FlowerbedMappers{
		mapper:       fMapper,
		imgMapper:    iMapper,
		sensorMapper: sMapper,
	}
}

func NewFlowerbedRepository(querier sqlc.Querier, mappers FlowerbedMappers) storage.FlowerbedRepository {
	return &FlowerbedRepository{
		querier:          querier,
		FlowerbedMappers: mappers,
	}
}

func (r *FlowerbedRepository) GetAll(ctx context.Context) ([]*entities.Flowerbed, error) {
	row, err := r.querier.GetAllFlowerbeds(ctx)
	if err != nil {
		return nil, err
	}

	mapped := r.mapper.FromSqlList(row)
	for _, f := range mapped {
		images, err := r.GetAllImagesByID(ctx, f.ID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get images")
		}
		f.Images = images

		sensor, err := r.GetSensorByFlowerbedID(ctx, f.ID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get sensor")
		}
		f.Sensor = sensor
	}

	return r.mapper.FromSqlList(row), nil
}

func (r *FlowerbedRepository) GetByID(ctx context.Context, id int32) (*entities.Flowerbed, error) {
	row, err := r.querier.GetFlowerbedByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return r.mapper.FromSql(row), nil
}

func (r *FlowerbedRepository) GetAllImagesByID(ctx context.Context, flowerbedID int32) ([]*entities.Image, error) {
	row, err := r.querier.GetAllImagesByFlowerbedID(ctx, flowerbedID)
	if err != nil {
		return nil, err
	}

	return r.imgMapper.FromSqlList(row), nil
}

func (r *FlowerbedRepository) GetSensorByFlowerbedID(ctx context.Context, flowerbedID int32) (*entities.Sensor, error) {
	row, err := r.querier.GetSensorByFlowerbedID(ctx, flowerbedID)
	if err != nil {
		return nil, err
	}

	return r.sensorMapper.FromSql(row), nil
}

func (r *FlowerbedRepository) Create(ctx context.Context, f *entities.Flowerbed) (*entities.Flowerbed, error) {
	params := sqlc.CreateFlowerbedParams{
		SensorID:       &f.Sensor.ID,
		Size:           f.Size,
		Description:    f.Description,
		NumberOfPlants: f.NumberOfPlants,
		MoistureLevel:  f.MoistureLevel,
		Region:         f.Region,
		Address:        f.Address,
		Latitude:       f.Latitude,
		Longitude:      f.Longitude,
	}

	row, err := r.querier.CreateFlowerbed(ctx, &params)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, row)
}

func (r *FlowerbedRepository) Update(ctx context.Context, f *entities.Flowerbed) (*entities.Flowerbed, error) {
	params := sqlc.UpdateFlowerbedParams{
		ID:             f.ID,
		SensorID:       &f.Sensor.ID,
		Size:           f.Size,
		Description:    f.Description,
		NumberOfPlants: f.NumberOfPlants,
		MoistureLevel:  f.MoistureLevel,
		Region:         f.Region,
		Address:        f.Address,
		Latitude:       f.Latitude,
		Longitude:      f.Longitude,
	}

	err := r.querier.UpdateFlowerbed(ctx, &params)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, f.ID)
}

func (r *FlowerbedRepository) Delete(ctx context.Context, id int32) error {
	return r.querier.DeleteFlowerbed(ctx, id)
}

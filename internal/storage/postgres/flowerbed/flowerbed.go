package flowerbed

import (
	"context"
	"log/slog"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/flowerbed/mapper"
	imgMapper "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/image/mapper"
	sensorMapper "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/sensor/mapper"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/pkg/errors"
)

type FlowerbedRepository struct {
	querier *sqlc.Queries
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

func NewFlowerbedRepository(querier *sqlc.Queries, mappers FlowerbedMappers) storage.FlowerbedRepository {
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
		if err := mapSensorAndImages(ctx, r, f); err != nil {
			return nil, err
		}
	}

	return mapped, nil
}

func (r *FlowerbedRepository) GetByID(ctx context.Context, id int32) (*entities.Flowerbed, error) {
	row, err := r.querier.GetFlowerbedByID(ctx, id)
	if err != nil {
		// TODO: Validate the error, it can also be a other error
		return nil, storage.ErrEntityNotFound
	}

	f := r.mapper.FromSql(row)
	if err := mapSensorAndImages(ctx, r, f); err != nil {
		return nil, err
	}

	return f, nil
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

func (r *FlowerbedRepository) Create(ctx context.Context, f *entities.CreateFlowerbed) (*entities.Flowerbed, error) {
	params := sqlc.CreateFlowerbedParams{
		SensorID:       f.SensorID,
		Size:           f.Size,
		Description:    f.Description,
		NumberOfPlants: f.NumberOfPlants,
		MoistureLevel:  f.MoistureLevel,
		Region:         f.Region,
		Address:        f.Address,
		Latitude:       f.Latitude,
		Longitude:      f.Longitude,
	}

	if f.SensorID != nil {
		_, err := r.querier.GetSensorByID(ctx, *f.SensorID)
		if err != nil {
			return nil, storage.ErrSensorNotFound
		}
	}

	for _, imgID := range f.ImageIDs {
		_, err := r.querier.GetImageByID(ctx, imgID)
		if err != nil {
			return nil, storage.ErrImageNotFound
		}
	}

	row, err := r.querier.CreateFlowerbed(ctx, &params)
	if err != nil {
		return nil, err
	}

	for _, imgID := range f.ImageIDs {
		params := sqlc.LinkFlowerbedImageParams{
			FlowerbedID: row,
			ImageID:     imgID,
		}
		if err = r.querier.LinkFlowerbedImage(ctx, &params); err != nil {
			return nil, err
		}
	}

	return r.GetByID(ctx, row)
}

func (r *FlowerbedRepository) Update(ctx context.Context, f *entities.UpdateFlowerbed) (*entities.Flowerbed, error) {
	prev, err := r.GetByID(ctx, f.ID)
	if err != nil {
		return nil, err
	}

	var sensorID *int32
	if f.SensorID != nil && prev.Sensor != nil {
		newSensorID := utils.CompareAndUpdate(prev.Sensor.ID, f.SensorID)
		sensorID = &newSensorID
		_, err = r.querier.GetSensorByID(ctx, newSensorID) // Check if sensor exists
		if err != nil {
			slog.Error("failed to get sensor by id. No sensor will be set", "error", err)
			sensorID = nil
		}
	} else if prev.Sensor != nil && f.SensorID == nil {
		sensorID = &prev.Sensor.ID
	} else if prev.Sensor == nil && f.SensorID != nil {
		sensorID = f.SensorID
	} else {
		sensorID = nil
	}

	for _, imgID := range f.ImageIDs {
		_, err := r.querier.GetImageByID(ctx, imgID) // Check if image exists
		if err != nil {
			return nil, storage.ErrImageNotFound
		}
	}

	params := sqlc.UpdateFlowerbedParams{
		ID:             f.ID,
		SensorID:       sensorID,
		Size:           utils.CompareAndUpdate(prev.Size, f.Size),
		Description:    utils.CompareAndUpdate(prev.Description, f.Description),
		NumberOfPlants: utils.CompareAndUpdate(prev.NumberOfPlants, f.NumberOfPlants),
		MoistureLevel:  utils.CompareAndUpdate(prev.MoistureLevel, f.MoistureLevel),
		Region:         utils.CompareAndUpdate(prev.Region, f.Region),
		Address:        utils.CompareAndUpdate(prev.Address, f.Address),
		Latitude:       utils.CompareAndUpdate(prev.Latitude, f.Latitude),
		Longitude:      utils.CompareAndUpdate(prev.Longitude, f.Longitude),
	}
	err = r.querier.UpdateFlowerbed(ctx, &params)
	if err != nil {
		return nil, err
	}

	if err = r.updateImages(ctx, prev, f); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, f.ID)
}

func (r *FlowerbedRepository) updateImages(ctx context.Context, prev *entities.Flowerbed, f *entities.UpdateFlowerbed) error {
	if f.ImageIDs == nil {
		return nil
	}

  // Unlink the images that are not in the new list
  for _, img := range prev.Images {
    found := false
    for _, newImgID := range f.ImageIDs {
      if img.ID == newImgID {
        found = true
        break
      }
    }

    if !found {
      args := sqlc.UnlinkFlowerbedImageParams{
        FlowerbedID: f.ID,
        ImageID:     img.ID,
      }
      if err := r.querier.UnlinkFlowerbedImage(ctx, &args); err != nil {
        return errors.Wrap(err, "failed to unlink image")
      }
    }
  }

  // Link the images that are in the new list
  for _, newImgID := range f.ImageIDs {
    found := false
    for _, img := range prev.Images {
      if img.ID == newImgID {
        found = true
        break
      }
    }

    if !found {
      args := sqlc.LinkFlowerbedImageParams{
        FlowerbedID: f.ID,
        ImageID:     newImgID,
      }
      if err := r.querier.LinkFlowerbedImage(ctx, &args); err != nil {
        return errors.Wrap(err, "failed to link image")
      }
    }
  }

	return nil
}

func (r *FlowerbedRepository) Delete(ctx context.Context, id int32) error {
	images, err := r.GetAllImagesByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to get images")
	}

	for _, img := range images {
		args := sqlc.UnlinkFlowerbedImageParams{
			FlowerbedID: id,
			ImageID:     img.ID,
		}
		if err = r.querier.UnlinkFlowerbedImage(ctx, &args); err != nil {
			return errors.Wrap(err, "failed to unlink image")
		}

		if err = r.querier.DeleteImage(ctx, img.ID); err != nil {
			return errors.Wrap(err, "failed to delete image")
		}
	}

	return r.querier.DeleteFlowerbed(ctx, id)
}

// Map sensor and images entity to domain flowerbed
func mapSensorAndImages(ctx context.Context, r *FlowerbedRepository, f *entities.Flowerbed) error {
	if err := mapImages(ctx, r, f); err != nil {
		return err
	}

	if err := mapSensor(ctx, r, f); err != nil {
		return err
	}

	return nil
}

func mapImages(ctx context.Context, r *FlowerbedRepository, f *entities.Flowerbed) error {
	images, err := r.GetAllImagesByID(ctx, f.ID)
	if err != nil {
		return errors.Wrap(err, "failed to get images")
	}
	f.Images = images
	return nil
}

func mapSensor(ctx context.Context, r *FlowerbedRepository, f *entities.Flowerbed) error {
	sensor, err := r.GetSensorByFlowerbedID(ctx, f.ID)
	if err != nil {
		return errors.Wrap(err, "failed to get sensor")
	}
	f.Sensor = sensor
	return nil
}

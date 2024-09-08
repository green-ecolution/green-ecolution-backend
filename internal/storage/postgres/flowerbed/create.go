package flowerbed

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

func defaultFlowerbed() *entities.Flowerbed {
	return &entities.Flowerbed{
		Sensor:         nil,
		Size:           0,
		Description:    "",
		NumberOfPlants: 0,
		MoistureLevel:  0,
		Region:         "",
		Address:        "",
		Latitude:       0,
		Longitude:      0,
		Images:         make([]*entities.Image, 0),
    Archived:       false,
	}
}

func (r *FlowerbedRepository) Create(ctx context.Context, fFn ...entities.EntityFunc[entities.Flowerbed]) (*entities.Flowerbed, error) {
	entity := defaultFlowerbed()
	for _, fn := range fFn {
		fn(entity)
	}

	id, err := r.createEntity(ctx, entity)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	entity.ID = *id

	return r.GetByID(ctx, *id)
}

func (r *FlowerbedRepository) createEntity(ctx context.Context, entity *entities.Flowerbed) (*int32, error) {
	args := sqlc.CreateFlowerbedParams{
		SensorID:       &entity.Sensor.ID,
		Size:           entity.Size,
		Description:    entity.Description,
		NumberOfPlants: entity.NumberOfPlants,
		MoistureLevel:  entity.MoistureLevel,
		Region:         entity.Region,
		Address:        entity.Address,
		Latitude:       entity.Latitude,
		Longitude:      entity.Longitude,
	}

	id, err := r.store.CreateFlowerbed(ctx, &args)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	return &id, nil
}

func (r *FlowerbedRepository) handleImages(ctx context.Context, flowerbedID int32, images []*entities.Image) error {
	for _, img := range images {
		err := r.linkImages(ctx, flowerbedID, img.ID)
		if err != nil {
			return err
		}
	}

  return nil
}

func (r *FlowerbedRepository) linkImages(ctx context.Context, flowerbedID int32, imgID int32) error {
	params := sqlc.LinkFlowerbedImageParams{
		FlowerbedID: flowerbedID,
		ImageID:     imgID,
	}
	return r.store.LinkFlowerbedImage(ctx, &params)
}

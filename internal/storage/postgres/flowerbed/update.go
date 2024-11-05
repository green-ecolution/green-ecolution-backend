package flowerbed

import (
	"context"
	"reflect"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/pkg/errors"
)

func (r *FlowerbedRepository) Update(ctx context.Context, id int32, fFn ...entities.EntityFunc[entities.Flowerbed]) (*entities.Flowerbed, error) {
	prev, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	original := *prev

	for _, fn := range fFn {
		fn(prev)
	}

	if reflect.DeepEqual(original, *prev) {
		return prev, nil
	}

	err = r.updateEntity(ctx, prev)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *FlowerbedRepository) UpdateWithImages(ctx context.Context, id int32, fFn ...entities.EntityFunc[entities.Flowerbed]) (*entities.Flowerbed, error) {	
	f, err := r.Update(ctx, id, fFn...)
	if err != nil {
		return nil, err
	}

	if err := r.updateImages(ctx, f); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *FlowerbedRepository) updateEntity(ctx context.Context, f *entities.Flowerbed) error {
	var sensorID *int32
	if f.Sensor != nil {
		sensorID = &f.Sensor.ID
	}

	args := sqlc.UpdateFlowerbedParams{
		ID:             f.ID,
		SensorID:       sensorID,
		Size:           f.Size,
		Description:    f.Description,
		NumberOfPlants: f.NumberOfPlants,
		MoistureLevel:  f.MoistureLevel,
		RegionID:       &f.Region.ID,
		Address:        f.Address,
		Latitude:       *f.Latitude,
		Longitude:      *f.Longitude,
	}

	return r.store.UpdateFlowerbed(ctx, &args)
}

func (r *FlowerbedRepository) updateImages(ctx context.Context, f *entities.Flowerbed) error {
	images, err := r.GetAllImagesByID(ctx, f.ID)
	if err != nil {
		return err
	}

	if (len(images) > 0) {
		if err := r.UnlinkAllImages(ctx, f.ID); err != nil {
			return err
		}
	}

	for _, img := range f.Images {
		if r.linkImages(ctx, f.ID, img.ID) != nil {
			return errors.New("error linking image")
		}
	}

	return nil
}

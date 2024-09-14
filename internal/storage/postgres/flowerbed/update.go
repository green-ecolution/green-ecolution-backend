package flowerbed

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/pkg/errors"
)

func (r *FlowerbedRepository) Update(ctx context.Context, id int32, fFn ...entities.EntityFunc[entities.Flowerbed]) (*entities.Flowerbed, error) {
	prev, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	for _, fn := range fFn {
		fn(prev)
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
	args := sqlc.UpdateFlowerbedParams{
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

	return r.store.UpdateFlowerbed(ctx, &args)
}

func (r *FlowerbedRepository) updateImages(ctx context.Context, f *entities.Flowerbed) error {
	if err := r.UnlinkAllImages(ctx, f.ID); err != nil {
		return err
	}

	for _, img := range f.Images {
		if r.linkImages(ctx, f.ID, img.ID) != nil {
			return errors.New("error linking image")
		}
	}

	return nil
}

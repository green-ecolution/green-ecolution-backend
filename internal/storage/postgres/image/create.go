package image

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

func defaultImage() *entities.Image {
	return &entities.Image{
		URL:      "",
		Filename: nil,
		MimeType: nil,
	}
}

func (r *ImageRepository) Create(ctx context.Context, iFn ...entities.EntityFunc[entities.Image]) (*entities.Image, error) {
	entity := defaultImage()
	for _, fn := range iFn {
		fn(entity)
	}

	id, err := r.createEntity(ctx, entity)
	if err != nil {
		return nil, err
	}

	entity.ID = *id
	return r.GetByID(ctx, *id)
}

func (r *ImageRepository) createEntity(ctx context.Context, image *entities.Image) (*int32, error) {
	args := sqlc.CreateImageParams{
		Url:      image.URL,
		Filename: image.Filename,
		MimeType: image.MimeType,
	}

	id, err := r.store.CreateImage(ctx, &args)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

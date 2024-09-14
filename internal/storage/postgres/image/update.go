package image

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

func (r *ImageRepository) Update(ctx context.Context, id int32, iFn ...entities.EntityFunc[entities.Image]) (*entities.Image, error) {
	entity, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	for _, fn := range iFn {
		fn(entity)
	}

	if err := r.updateEntity(ctx, entity); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, entity.ID)
}

func (r *ImageRepository) updateEntity(ctx context.Context, image *entities.Image) error {
	params := sqlc.UpdateImageParams{
		ID:       image.ID,
		Url:      image.URL,
		Filename: image.Filename,
		MimeType: image.MimeType,
	}

	return r.store.UpdateImage(ctx, &params)
}

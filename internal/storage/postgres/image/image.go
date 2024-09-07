package image

import (
	"context"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	. "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
)

type ImageRepository struct {
	Store *Store
	ImageRepositoryMappers
}

type ImageRepositoryMappers struct {
	mapper mapper.InternalImageRepoMapper
}

func NewImageRepositoryMappers(iMapper mapper.InternalImageRepoMapper) ImageRepositoryMappers {
	return ImageRepositoryMappers{
		mapper: iMapper,
	}
}

func NewImageRepository(store *Store, mappers ImageRepositoryMappers) storage.ImageRepository {
	store.SetEntityType(Image)
	return &ImageRepository{
		Store:                  store,
		ImageRepositoryMappers: mappers,
	}
}

func (r *ImageRepository) GetAll(ctx context.Context) ([]*entities.Image, error) {
	rows, err := r.Store.GetAllImages(ctx)
	if err != nil {
		return nil, r.Store.HandleError(err)
	}

	return r.mapper.FromSqlList(rows), nil
}

func (r *ImageRepository) GetByID(ctx context.Context, id int32) (*entities.Image, error) {
	row, err := r.Store.GetImageByID(ctx, id)
	if err != nil {
		return nil, r.Store.HandleError(err)
	}

	return r.mapper.FromSql(row), nil
}

func (r *ImageRepository) Create(ctx context.Context, image *entities.CreateImage) (*entities.Image, error) {
	params := &sqlc.CreateImageParams{
		Url:      image.URL,
		Filename: image.Filename,
		MimeType: image.MimeType,
	}
	id, err := r.Store.CreateImage(ctx, params)
	if err != nil {
		return nil, r.Store.HandleError(err)
	}

	return r.GetByID(ctx, id)
}

func (r *ImageRepository) Update(ctx context.Context, image *entities.UpdateImage) (*entities.Image, error) {
	prev, err := r.GetByID(ctx, image.ID)
	if err != nil {
		return nil, r.Store.HandleError(err)
	}

	if !hasChanges(prev, image) {
		return prev, nil
	}

	if image.URL == nil {
		image.URL = &prev.URL
	}

	params := &sqlc.UpdateImageParams{
		ID:       image.ID,
		Url:      *image.URL,
		Filename: utils.CompareAndUpdate(prev.Filename, &image.Filename),
		MimeType: utils.CompareAndUpdate(prev.MimeType, &image.MimeType),
	}

	if err := r.Store.UpdateImage(ctx, params); err != nil {
		return nil, r.Store.HandleError(err)
	}

	return r.GetByID(ctx, image.ID)
}

func (r *ImageRepository) Delete(ctx context.Context, id int32) error {
	return r.Store.DeleteImage(ctx, id)
}

func hasChanges(p *entities.Image, u *entities.UpdateImage) bool {
	var url string
	if u.URL != nil {
		url = *u.URL
	} else {
		url = p.URL
	}

	return (u.Filename != nil && *u.Filename != *p.Filename) ||
		(u.MimeType != nil && *u.MimeType != *p.MimeType) ||
		url != p.URL
}

package image

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/image/mapper"
)

type ImageRepository struct {
	querier sqlc.Querier
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

func NewImageRepository(querier sqlc.Querier, mappers ImageRepositoryMappers) storage.ImageRepository {
	return &ImageRepository{
		querier:                querier,
		ImageRepositoryMappers: mappers,
	}
}

func (r *ImageRepository) GetAll(ctx context.Context) ([]*entities.Image, error) {
	rows, err := r.querier.GetAllImages(ctx)
	if err != nil {
		return nil, err
	}

	return r.mapper.FromSqlList(rows), nil
}

func (r *ImageRepository) GetByID(ctx context.Context, id int32) (*entities.Image, error) {
	row, err := r.querier.GetImageByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return r.mapper.FromSql(row), nil
}

func (r *ImageRepository) Create(ctx context.Context, image *entities.Image) (*entities.Image, error) {
	params := &sqlc.CreateImageParams{
		Url:      image.URL,
		Filename: image.Filename,
		MimeType: image.MimeType,
	}
	id, err := r.querier.CreateImage(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *ImageRepository) Update(ctx context.Context, image *entities.Image) (*entities.Image, error) {
	params := &sqlc.UpdateImageParams{
		ID:       image.ID,
		Url:      image.URL,
		Filename: image.Filename,
		MimeType: image.MimeType,
	}

	if err := r.querier.UpdateImage(ctx, params); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, image.ID)
}

func (r *ImageRepository) Delete(ctx context.Context, id int32) error {
	return r.querier.DeleteImage(ctx, id)
}

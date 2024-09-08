package image

import (
	"context"
	"log/slog"

	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
)

type ImageRepository struct {
	store *store.Store
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

func NewImageRepository(s *store.Store, mappers ImageRepositoryMappers) storage.ImageRepository {
	s.SetEntityType(store.Image)
	return &ImageRepository{
		store:                  s,
		ImageRepositoryMappers: mappers,
	}
}

func WithURL(url string) entities.EntityFunc[entities.Image] {
	return func(i *entities.Image) {
		slog.Debug("updating url", "url", url)
		i.URL = url
	}
}

func WithFilename(filename *string) entities.EntityFunc[entities.Image] {
	return func(i *entities.Image) {
		slog.Debug("updating filename", "filename", filename)
		i.Filename = filename
	}
}

func WithMimeType(mimeType *string) entities.EntityFunc[entities.Image] {
	return func(i *entities.Image) {
		slog.Debug("updating mime type", "mime type", mimeType)
		i.MimeType = mimeType
	}
}

func (r *ImageRepository) Delete(ctx context.Context, id int32) error {
	return r.store.DeleteImage(ctx, id)
}

package tree

import (
	"context"
	imgMapper "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
	"github.com/pkg/errors"
)

type TreeRepository struct {
	store *store.Store
	TreeMappers
}

type TreeMappers struct {
	mapper   imgMapper.InternalTreeRepoMapper
	iMapper  imgMapper.InternalImageRepoMapper
	sMapper  imgMapper.InternalSensorRepoMapper
	tcMapper imgMapper.InternalTreeClusterRepoMapper
}

func NewTreeRepositoryMappers(
	tMapper imgMapper.InternalTreeRepoMapper,
	iMapper imgMapper.InternalImageRepoMapper,
	sMapper imgMapper.InternalSensorRepoMapper,
	tcMapper imgMapper.InternalTreeClusterRepoMapper,
) TreeMappers {
	return TreeMappers{
		mapper:   tMapper,
		iMapper:  iMapper,
		sMapper:  sMapper,
		tcMapper: tcMapper,
	}
}

func NewTreeRepository(s *store.Store, mappers TreeMappers) storage.TreeRepository {
	s.SetEntityType(store.Tree)
	return &TreeRepository{
		store:       s,
		TreeMappers: mappers,
	}
}

func WithSpecies(species string) entities.EntityFunc[entities.Tree] {
	return func(t *entities.Tree) {
		t.Species = species
	}
}

func WithAge(age int32) entities.EntityFunc[entities.Tree] {
	return func(t *entities.Tree) {
		t.Age = age
	}
}

func WithHeightAboveSeaLevel(height float64) entities.EntityFunc[entities.Tree] {
	return func(t *entities.Tree) {
		t.HeightAboveSeaLevel = height
	}
}

func WithSensor(sensor *entities.Sensor) entities.EntityFunc[entities.Tree] {
	return func(t *entities.Tree) {
		t.Sensor = sensor
	}
}

func WithPlantingYear(year int32) entities.EntityFunc[entities.Tree] {
	return func(t *entities.Tree) {
		t.PlantingYear = year
	}
}

func WithLatitude(lat float64) entities.EntityFunc[entities.Tree] {
	return func(t *entities.Tree) {
		t.Latitude = lat
	}
}

func WithLongitude(long float64) entities.EntityFunc[entities.Tree] {
	return func(t *entities.Tree) {
		t.Longitude = long
	}
}

func WithTreeNumber(number int32) entities.EntityFunc[entities.Tree] {
	return func(t *entities.Tree) {
		t.Number = number
	}
}

func WithTreeCluster(treeCluster *entities.TreeCluster) entities.EntityFunc[entities.Tree] {
	return func(t *entities.Tree) {
		t.TreeCluster = treeCluster
	}
}

func WithImages(images []*entities.Image) entities.EntityFunc[entities.Tree] {
	return func(t *entities.Tree) {
		t.Images = images
	}
}

func (r *TreeRepository) Delete(ctx context.Context, id int32) error {
	images, err := r.GetAllImagesByID(ctx, id)
	if err != nil {
		return r.store.HandleError(errors.Wrap(err, "failed to get images"))
	}

	for _, img := range images {
		args := sqlc.UnlinkTreeImageParams{
			TreeID:  id,
			ImageID: img.ID,
		}
		if err = r.store.UnlinkTreeImage(ctx, &args); err != nil {
			return r.store.HandleError(errors.Wrap(err, "failed to unlink image"))
		}

		if err = r.store.DeleteImage(ctx, img.ID); err != nil {
			return r.store.HandleError(errors.Wrap(err, "failed to delete image"))
		}
	}
	return r.store.DeleteTree(ctx, id)
}

func (r *TreeRepository) DeleteAndUnlinkImages(ctx context.Context, id int32) error {
	if err := r.UnlinkAllImages(ctx, id); err != nil {
		return r.store.HandleError(errors.Wrap(err, "failed to unlink images"))
	}

	return r.Delete(ctx, id)
}

func (r *TreeRepository) UnlinkImage(ctx context.Context, treeID, imageID int32) error {
	args := sqlc.UnlinkTreeImageParams{
		TreeID:  treeID,
		ImageID: imageID,
	}
	return r.store.UnlinkTreeImage(ctx, &args)
}

func (r *TreeRepository) UnlinkAllImages(ctx context.Context, treeID int32) error {
	return r.store.UnlinkAllTreeImages(ctx, treeID)
}

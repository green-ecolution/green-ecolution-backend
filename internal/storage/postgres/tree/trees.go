package tree

import (
	"context"
	"errors"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	imgMapper "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
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

func WithProvider(provider string) entities.EntityFunc[entities.Tree] {
	return func(t *entities.Tree) {
		t.Provider = provider
	}
}

func WithAdditionalInfo(additionalInfo map[string]interface{}) entities.EntityFunc[entities.Tree] {
	return func(t *entities.Tree) {
		t.AdditionalInfo = additionalInfo
	}
}

func WithDescription(description string) entities.EntityFunc[entities.Tree] {
	return func(t *entities.Tree) {
		t.Description = description
	}
}

func WithReadonly(readonly bool) entities.EntityFunc[entities.Tree] {
	return func(t *entities.Tree) {
		t.Readonly = readonly
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

func WithNumber(number string) entities.EntityFunc[entities.Tree] {
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

func WithWateringStatus(wateringStatus entities.WateringStatus) entities.EntityFunc[entities.Tree] {
	return func(tc *entities.Tree) {
		tc.WateringStatus = wateringStatus
	}
}

func (r *TreeRepository) Delete(ctx context.Context, id int32) error {
	log := logger.GetLogger(ctx)
	images, err := r.GetAllImagesByID(ctx, id)
	if err != nil {
		return err
	}

	for _, img := range images {
		args := sqlc.UnlinkTreeImageParams{
			TreeID:  id,
			ImageID: img.ID,
		}
		_, err = r.store.UnlinkTreeImage(ctx, &args)
		if err != nil {
			log.Debug("failed to unlink image of tree in db", "error", err, "tree_id", id, "image_id", img.ID)
			return err
		}

		if err := r.store.DeleteImage(ctx, img.ID); err != nil {
			log.Debug("failed to delete image of tree in db", "error", err, "tree_id", id, "image_id", img.ID)
			return err
		}
	}

	_, err = r.store.DeleteTree(ctx, id)
	if err != nil {
		log.Debug("failed to delete tree in db", "error", err, "tree_id", id)
		return err
	}

	log.Debug("tree entity deleted successfully in db", "tree_id", id)
	return nil
}

func (r *TreeRepository) DeleteAndUnlinkImages(ctx context.Context, id int32) error {
	if err := r.UnlinkAllImages(ctx, id); err != nil {
		return err
	}

	return r.Delete(ctx, id)
}

func (r *TreeRepository) UnlinkImage(ctx context.Context, treeID, imageID int32) error {
	args := sqlc.UnlinkTreeImageParams{
		TreeID:  treeID,
		ImageID: imageID,
	}
	_, err := r.store.UnlinkTreeImage(ctx, &args)
	return err
}

func (r *TreeRepository) UnlinkAllImages(ctx context.Context, treeID int32) error {
	return r.store.UnlinkAllTreeImages(ctx, treeID)
}

func (r *TreeRepository) UnlinkTreeClusterID(ctx context.Context, treeClusterID int32) error {
	log := logger.GetLogger(ctx)

	_, err := r.store.GetTreeClusterByID(ctx, treeClusterID)
	if err != nil {
		return err
	}
	unlinkTreeIDs, err := r.store.UnlinkTreeClusterID(ctx, &treeClusterID)
	if err != nil {
		log.Error("failed to unlink tree cluster from trees", "error", err, "cluster_id", treeClusterID)
	}

	log.Info("unlink trees from following tree cluster", "cluster_id", treeClusterID, "unlinked_trees", unlinkTreeIDs)

	return nil
}

func (r *TreeRepository) UnlinkSensorID(ctx context.Context, sensorID string) error {
	if sensorID == "" {
		return errors.New("sensorID cannot be empty")
	}
	return r.store.UnlinkSensorIDFromTrees(ctx, &sensorID)
}

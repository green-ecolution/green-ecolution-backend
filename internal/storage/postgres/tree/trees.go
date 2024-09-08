package tree

import (
	"context"
	imgMapper "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper"
	"log/slog"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/jackc/pgx/v5"
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

func (r *TreeRepository) GetAll(ctx context.Context) ([]*entities.Tree, error) {
	rows, err := r.store.GetAllTrees(ctx)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	t := r.mapper.FromSqlList(rows)
	for _, tree := range t {
		if err := r.mapFields(ctx, tree); err != nil {
			return nil, r.store.HandleError(err)
		}
	}

	return t, nil
}

func (r *TreeRepository) GetByID(ctx context.Context, id int32) (*entities.Tree, error) {
	row, err := r.store.GetTreeByID(ctx, id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	t := r.mapper.FromSql(row)
	if err := r.mapFields(ctx, t); err != nil {
		return nil, r.store.HandleError(err)
	}

	return t, nil
}

func (r *TreeRepository) GetByTreeClusterID(ctx context.Context, id int32) ([]*entities.Tree, error) {
	rows, err := r.store.GetTreesByTreeClusterID(ctx, &id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	t := r.mapper.FromSqlList(rows)
	for _, tree := range t {
		if err := r.mapFields(ctx, tree); err != nil {
			return nil, r.store.HandleError(err)
		}
	}

	return t, nil
}

func (r *TreeRepository) GetAllImagesByID(ctx context.Context, id int32) ([]*entities.Image, error) {
	rows, err := r.store.GetAllImagesByTreeID(ctx, id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	return r.iMapper.FromSqlList(rows), nil
}

func (r *TreeRepository) Create(ctx context.Context, tree *entities.CreateTree) (*entities.Tree, error) {
	entity := sqlc.CreateTreeParams{
		TreeClusterID:       &tree.TreeClusterID,
		Species:             tree.Species,
		Age:                 tree.Age,
		HeightAboveSeaLevel: tree.HeightAboveSeaLevel,
		SensorID:            tree.SensorID,
		PlantingYear:        tree.PlantingYear,
		Latitude:            tree.Latitude,
		Longitude:           tree.Longitude,
	}

	id, err := r.store.CreateTree(ctx, &entity)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	// Link images to tree
	for _, img := range tree.ImageIDs {
		if err := r.store.CheckImageExists(ctx, img); err != nil {
			if errors.Is(err, storage.ErrImageNotFound) {
				slog.Error("failed to get image by id. No image will be set", "imageID", img, "error", err)
				continue
			} else {
				return nil, err
			}
		}

		params := sqlc.LinkTreeImageParams{
			TreeID:  id,
			ImageID: *img,
		}
		if err = r.store.LinkTreeImage(ctx, &params); err != nil {
			return nil, r.store.HandleError(err)
		}
	}

	return r.GetByID(ctx, id)
}

func (r *TreeRepository) Update(ctx context.Context, tree *entities.UpdateTree) (*entities.Tree, error) {
	prev, err := r.GetByID(ctx, tree.ID)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	// Check if sensor exists and update sensorID
	var sensorID *int32
	if tree.SensorID != nil && prev.Sensor != nil {
		newSensorID := utils.CompareAndUpdate(prev.Sensor.ID, tree.SensorID)
		sensorID = &newSensorID
		_, err = r.store.GetSensorByID(ctx, newSensorID) // Check if sensor exists
		if err != nil {
			if errors.Is(err, storage.ErrSensorNotFound) {
				slog.Error("failed to get sensor by id. No sensor will be set", "error", err)
				sensorID = nil
			} else {
				return nil, r.store.HandleError(err)
			}
		}
	} else if prev.Sensor != nil && tree.SensorID == nil {
		sensorID = &prev.Sensor.ID
	} else if prev.Sensor == nil && tree.SensorID != nil {
		sensorID = tree.SensorID
	} else {
		sensorID = nil
	}

	entity := sqlc.UpdateTreeParams{
		ID:                  tree.ID,
		Species:             utils.CompareAndUpdate(prev.Species, tree.Species),
		Age:                 utils.CompareAndUpdate(prev.Age, tree.Age),
		HeightAboveSeaLevel: utils.CompareAndUpdate(prev.HeightAboveSeaLevel, tree.HeightAboveSeaLevel),
		SensorID:            sensorID,
		PlantingYear:        utils.CompareAndUpdate(prev.PlantingYear, tree.PlantingYear),
		Latitude:            utils.CompareAndUpdate(prev.Latitude, tree.Latitude),
		Longitude:           utils.CompareAndUpdate(prev.Longitude, tree.Longitude),
		TreeNumber:          utils.CompareAndUpdate(prev.Number, tree.Number),
		TreeClusterID:       utils.P(utils.CompareAndUpdate(prev.TreeCluster.ID, tree.TreeClusterID)),
	}

	err = r.store.UpdateTree(ctx, &entity)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	// Update Images
	if err := r.updateImages(ctx, prev, tree); err != nil {
		return nil, r.store.HandleError(err)
	}

	return r.GetByID(ctx, tree.ID)
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

func (r *TreeRepository) updateImages(ctx context.Context, prev *entities.Tree, f *entities.UpdateTree) error {
	if f.ImageIDs == nil {
		return nil
	}

	// Unlink the images that are not in the new list
	for _, img := range prev.Images {
		found := false
		for _, newImgID := range f.ImageIDs {
			if img.ID == *newImgID {
				found = true
				break
			}
		}

		if !found {
			args := sqlc.UnlinkFlowerbedImageParams{
				FlowerbedID: f.ID,
				ImageID:     img.ID,
			}
			if err := r.store.UnlinkFlowerbedImage(ctx, &args); err != nil {
				return r.store.HandleError(errors.Wrap(err, "failed to unlink image"))
			}
		}
	}

	// Link the images that are in the new list
	for _, newImgID := range f.ImageIDs {
		found := false
		for _, img := range prev.Images {
			if img.ID == *newImgID {
				found = true
				break
			}
		}

		if !found {
			args := sqlc.LinkFlowerbedImageParams{
				FlowerbedID: f.ID,
				ImageID:     *newImgID,
			}
			if err := r.store.LinkFlowerbedImage(ctx, &args); err != nil {
				return r.store.HandleError(errors.Wrap(err, "failed to unlink image"))
			}
		}
	}

	return nil
}

func (r *TreeRepository) GetSensorByTreeID(ctx context.Context, flowerbedID int32) (*entities.Sensor, error) {
	row, err := r.store.GetSensorByTreeID(ctx, flowerbedID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrSensorNotFound
		} else {
			return nil, r.store.HandleError(err)
		}
	}

	return r.sMapper.FromSql(row), nil
}

func (r *TreeRepository) GetTreeClusterByTreeID(ctx context.Context, treeID int32) (*entities.TreeCluster, error) {
	row, err := r.store.GetTreeClusterByTreeID(ctx, treeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrTreeClusterNotFound
		} else {
			return nil, r.store.HandleError(err)
		}
	}

	return r.tcMapper.FromSql(row), nil
}

// Map sensor, images and tree cluster entity to domain flowerbed
func (r *TreeRepository) mapFields(ctx context.Context, t *entities.Tree) error {
	if err := mapImages(ctx, r, t); err != nil {
		return r.store.HandleError(err)
	}

	if err := mapSensor(ctx, r, t); err != nil {
		return r.store.HandleError(err)
	}

	if err := mapTreeCluster(ctx, r, t); err != nil {
		return r.store.HandleError(err)
	}

	return nil
}

func mapImages(ctx context.Context, r *TreeRepository, t *entities.Tree) error {
	images, err := r.GetAllImagesByID(ctx, t.ID)
	if err != nil {
		return r.store.HandleError(err)
	}
	t.Images = images
	return nil
}

func mapSensor(ctx context.Context, r *TreeRepository, t *entities.Tree) error {
	sensor, err := r.GetSensorByTreeID(ctx, t.ID)
	if err != nil {
		if errors.Is(err, storage.ErrSensorNotFound) {
			// If sensor is not found, set sensor to nil
			t.Sensor = nil
			return nil
		}
		return r.store.HandleError(err)
	}
	t.Sensor = sensor
	return nil
}

func mapTreeCluster(ctx context.Context, r *TreeRepository, t *entities.Tree) error {
	treeCluster, err := r.GetTreeClusterByTreeID(ctx, t.ID)
	if err != nil {
		return r.store.HandleError(err)
	}
	t.TreeCluster = treeCluster
	return nil
}

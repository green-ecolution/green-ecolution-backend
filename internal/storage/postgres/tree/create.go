package tree

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

func defaultTree() entities.Tree {
	return entities.Tree{
		TreeCluster:    nil,
		Species:        "",
		Number:         "",
		Readonly:       false,
		Sensor:         nil,
		PlantingYear:   0,
		Latitude:       0,
		Longitude:      0,
		Images:         nil,
		WateringStatus: entities.WateringStatusUnknown,
		Description:    "",
	}
}

func (r *TreeRepository) Create(ctx context.Context, tFn ...entities.EntityFunc[entities.Tree]) (*entities.Tree, error) {
	entity := defaultTree()
	for _, fn := range tFn {
		fn(&entity)
	}

	id, err := r.createEntity(ctx, &entity)
	if err != nil {
		return nil, r.store.HandleError(err)
	}
	entity.ID = id

	return r.GetByID(ctx, id)
}

func (r *TreeRepository) CreateAndLinkImages(ctx context.Context, tFn ...entities.EntityFunc[entities.Tree]) (*entities.Tree, error) {
	createdEntity, err := r.Create(ctx, tFn...)
	if err != nil {
		return nil, err
	}

	entity := defaultTree()
	for _, fn := range tFn {
		fn(&entity)
	}

	if entity.Images != nil {
		if err := r.handleImages(ctx, createdEntity.ID, entity.Images); err != nil {
			return nil, err
		}
		linkedImages, err := r.GetAllImagesByID(ctx, createdEntity.ID)
		if err != nil {
			return nil, err
		}
		createdEntity.Images = linkedImages
	}

	return createdEntity, nil
}

func (r *TreeRepository) createEntity(ctx context.Context, entity *entities.Tree) (int32, error) {
	var treeClusterID *int32
	if entity.TreeCluster != nil {
		treeClusterID = &entity.TreeCluster.ID
	}

	var sensorID *int32
	if entity.Sensor != nil {
		sensorID = &entity.Sensor.ID
	}

	args := sqlc.CreateTreeParams{
		TreeClusterID:  treeClusterID,
		Species:        entity.Species,
		Readonly:       entity.Readonly,
		SensorID:       sensorID,
		PlantingYear:   entity.PlantingYear,
		Latitude:       entity.Latitude,
		Longitude:      entity.Longitude,
		WateringStatus: sqlc.WateringStatus(entity.WateringStatus),
		Description:    &entity.Description,
		TreeNumber:     entity.Number,
	}

	id, err := r.store.CreateTree(ctx, &args)
	if err != nil {
		return -1, err
	}

	if err := r.store.SetTreeLocation(ctx, &sqlc.SetTreeLocationParams{
		ID:        id,
		Latitude:  entity.Latitude,
		Longitude: entity.Longitude,
	}); err != nil {
		return -1, err
	}

	return id, nil
}

func (r *TreeRepository) handleImages(ctx context.Context, treeID int32, images []*entities.Image) error {
	for _, img := range images {
		err := r.linkImages(ctx, treeID, img.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *TreeRepository) linkImages(ctx context.Context, treeID, imgID int32) error {
	params := sqlc.LinkTreeImageParams{
		TreeID:  treeID,
		ImageID: imgID,
	}

	return r.store.LinkTreeImage(ctx, &params)
}

package tree

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

func defaultTree() entities.Tree {
	return entities.Tree{
		TreeCluster:         nil,
		Species:             "",
		Age:                 0,
		HeightAboveSeaLevel: 0,
		Sensor:              nil,
		PlantingYear:        0,
		Latitude:            0,
		Longitude:           0,
		Images:              nil,
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
	entity, err := r.Create(ctx, tFn...)
	if err != nil {
		return nil, err
	}

	if err := r.handleImages(ctx, entity.ID, entity.Images); err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *TreeRepository) createEntity(ctx context.Context, entity *entities.Tree) (int32, error) {
	args := sqlc.CreateTreeParams{
		TreeClusterID:       &entity.TreeCluster.ID,
		Species:             entity.Species,
		Age:                 entity.Age,
		HeightAboveSeaLevel: entity.HeightAboveSeaLevel,
		SensorID:            &entity.Sensor.ID,
		PlantingYear:        entity.PlantingYear,
		Latitude:            entity.Latitude,
		Longitude:           entity.Longitude,
	}

	return r.store.CreateTree(ctx, &args)
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

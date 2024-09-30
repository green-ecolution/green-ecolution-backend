package tree

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/pkg/errors"
)

func (r *TreeRepository) Update(ctx context.Context, id int32, tFn ...entities.EntityFunc[entities.Tree]) (*entities.Tree, error) {
	entity, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	for _, fn := range tFn {
		fn(entity)
	}

	err = r.updateEntity(ctx, entity)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *TreeRepository) UpdateWithImages(ctx context.Context, id int32, tFn ...entities.EntityFunc[entities.Tree]) (*entities.Tree, error) {
	t, err := r.Update(ctx, id, tFn...)
	if err != nil {
		return nil, err
	}

	if err := r.updateImages(ctx, t); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *TreeRepository) UpdateTreeClusterID(ctx context.Context, treeIDs []int32, treeClusterID *int32) error {
	args := &sqlc.UpdateTreeClusterIDParams{
		Column1:       treeIDs,
		TreeClusterID: treeClusterID,
	}

	return r.store.UpdateTreeClusterID(ctx, args)
}

func (r *TreeRepository) updateEntity(ctx context.Context, t *entities.Tree) error {
	args := sqlc.UpdateTreeParams{
		ID:                  t.ID,
		Species:             t.Species,
		Age:                 t.Age,
		HeightAboveSeaLevel: t.HeightAboveSeaLevel,
		SensorID:            &t.Sensor.ID,
		PlantingYear:        t.PlantingYear,
		Latitude:            t.Latitude,
		Longitude:           t.Longitude,
		TreeNumber:          t.Number,
		TreeClusterID:       &t.TreeCluster.ID,
	}

	return r.store.UpdateTree(ctx, &args)
}

func (r *TreeRepository) updateImages(ctx context.Context, tree *entities.Tree) error {
	if err := r.UnlinkAllImages(ctx, tree.ID); err != nil {
		return err
	}

	for _, img := range tree.Images {
		if r.linkImages(ctx, tree.ID, img.ID) != nil {
			return errors.New("error linking image")
		}
	}

	return nil
}

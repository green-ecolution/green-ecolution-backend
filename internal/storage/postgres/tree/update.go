package tree

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/pkg/errors"
)

func (r *TreeRepository) Update(ctx context.Context, id int32, tFn ...entities.EntityFunc[entities.Tree]) (*entities.Tree, error) {
	log := logger.GetLogger(ctx)
	entity, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	for _, fn := range tFn {
		fn(entity)
	}

	err = r.updateEntity(ctx, entity)
	if err != nil {
		log.Error("failed to update tree entity in db", "error", err, "tree_id", id)
		return nil, err
	}

	log.Debug("tree entity updated successfully in db", "tree_id", id)
	return r.GetByID(ctx, id)
}

func (r *TreeRepository) UpdateWithImages(ctx context.Context, id int32, tFn ...entities.EntityFunc[entities.Tree]) (*entities.Tree, error) {
	t, err := r.Update(ctx, id, tFn...)
	if err != nil {
		return nil, err
	}

	entity := defaultTree()
	for _, fn := range tFn {
		fn(&entity)
	}

	if len(entity.Images) > 0 {
		if t.Images == nil {
			t.Images = entity.Images
		} else {
			t.Images = append(t.Images, entity.Images...)
		}
		if err := r.updateImages(ctx, t); err != nil {
			return nil, err
		}
	}
	return r.GetByID(ctx, id)
}

func (r *TreeRepository) updateEntity(ctx context.Context, t *entities.Tree) error {
	log := logger.GetLogger(ctx)
	additionalInfo, err := utils.MapAdditionalInfoToByte(t.AdditionalInfo)
	if err != nil {
		log.Debug("failed to marshal additional informations to byte array", "error", err, "additional_info", t.AdditionalInfo)
		return err
	}

	var treeClusterID *int32
	if t.TreeCluster != nil {
		treeClusterID = &t.TreeCluster.ID
	}

	var sensorID *string
	if t.Sensor != nil {
		sensorID = &t.Sensor.ID

		if err := r.store.UnlinkSensorIDFromTrees(ctx, sensorID); err != nil {
			return err
		}
	}

	args := sqlc.UpdateTreeParams{
		ID:                     t.ID,
		Species:                t.Species,
		Readonly:               t.Readonly,
		PlantingYear:           t.PlantingYear,
		Number:                 t.Number,
		SensorID:               sensorID,
		TreeClusterID:          treeClusterID,
		WateringStatus:         sqlc.WateringStatus(t.WateringStatus),
		Description:            &t.Description,
		Provider:               &t.Provider,
		AdditionalInformations: additionalInfo,
	}

	if err := r.store.SetTreeLocation(ctx, &sqlc.SetTreeLocationParams{
		ID:        t.ID,
		Latitude:  t.Latitude,
		Longitude: t.Longitude,
	}); err != nil {
		return err
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

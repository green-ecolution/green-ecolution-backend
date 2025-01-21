package tree

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/pkg/errors"
)

func (r *TreeRepository) Update(ctx context.Context, id int32, updateFn func(*entities.Tree) (bool, error)) (*entities.Tree, error) {
	return r.updateWithTx(ctx, id, updateFn, nil)
}

func (r *TreeRepository) UpdateWithImages(ctx context.Context, id int32, updateFn func(*entities.Tree) (bool, error)) (*entities.Tree, error) {
	return r.updateWithTx(ctx, id, updateFn, func(ctx context.Context, tree *entities.Tree, updatedEntity *entities.Tree) error {
		if len(tree.Images) > 0 {
			if updatedEntity.Images == nil {
				updatedEntity.Images = tree.Images
			} else {
				updatedEntity.Images = append(updatedEntity.Images, tree.Images...)
			}
			if err := r.updateImages(ctx, updatedEntity); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *TreeRepository) updateWithTx(
	ctx context.Context,
	id int32,
	updateFn func(*entities.Tree) (bool, error),
	afterUpdateFn func(ctx context.Context, tree *entities.Tree, updatedEntity *entities.Tree) error) (*entities.Tree, error) {
	log := logger.GetLogger(ctx)
	var updatedTree *entities.Tree

	err := r.store.WithTx(ctx, func(s *store.Store) error {
		oldStore := r.store
		defer func() {
			r.store = oldStore
		}()
		r.store = s

		tree, err := r.GetByID(ctx, id)
		if err != nil {
			log.Error("failed to get tree entity from db", "error", err, "tree_id", id)
			return err
		}

		if updateFn == nil {
			return errors.New("updateFn is nil")
		}

		updated, err := updateFn(tree)
		if err != nil {
			return err
		}

		if !updated {
			updatedTree = tree
			return nil
		}

		if err := r.updateEntity(ctx, tree); err != nil {
			log.Error("failed to update tree entity in db", "error", err, "tree_id", id)
			return err
		}

		updatedTree, err = r.GetByID(ctx, id)
		if err != nil {
			log.Error("failed to get updated tree entity from db", "error", err, "tree_id", id)
			return err
		}

		if afterUpdateFn != nil {
			if err := afterUpdateFn(ctx, tree, updatedTree); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	log.Debug("tree entity updated successfully in db", "tree_id", id)
	return updatedTree, nil
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

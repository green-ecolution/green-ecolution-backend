package tree

import (
	"context"
	"errors"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
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

func (r *TreeRepository) Create(ctx context.Context, createFn func(*entities.Tree) (bool, error)) (*entities.Tree, error) {
	if createFn == nil {
		return nil, errors.New("createFn is nil")
	}

	var createdTr *entities.Tree
	err := r.store.WithTx(ctx, func(s *store.Store) error {
		oldStore := r.store
		defer func() {
			r.store = oldStore
		}()
		r.store = s

		entity := defaultTree()
		created, err := createFn(&entity)
		if err != nil {
			return err
		}

		if !created {
			return nil
		}

		if err := r.validateTreeEntity(&entity); err != nil {
			return err
		}

		id, err := r.createEntity(ctx, &entity)
		if err != nil {
			return err
		}

		createdTr, err = r.GetByID(ctx, id)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdTr, nil
}

func (r *TreeRepository) CreateAndLinkImages(ctx context.Context, createFn func(*entities.Tree) (bool, error)) (*entities.Tree, error) {
	if createFn == nil {
		return nil, errors.New("createFn is nil")
	}

	var createdTr *entities.Tree
	err := r.store.WithTx(ctx, func(s *store.Store) error {
		oldStore := r.store
		defer func() {
			r.store = oldStore
		}()
		r.store = s

		entity := defaultTree()
		created, err := createFn(&entity)
		if err != nil {
			return err
		}
		
		if !created {
			return nil
		}

		if err := r.validateTreeEntity(&entity); err != nil {
			return err
		}

		id, err := r.createEntity(ctx, &entity)
		if err != nil {
			return err
		}

		createdTr, err = r.GetByID(ctx, id)
		if err != nil {
			return err
		}

		if entity.Images != nil {
			if err := r.handleImages(ctx, id, entity.Images); err != nil {
				return err
			}

			linkedImages, err := r.GetAllImagesByID(ctx, id)
			if err != nil {
				return err
			}
			createdTr.Images = linkedImages
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdTr, nil
}

func (r *TreeRepository) transaction(
	ctx context.Context,
	createFn func(*entities.Tree) (bool, error),
	handler func(context.Context, *entities.Tree, string) error,
) (*entities.Tree, error) {
	if createFn == nil {
		return nil, errors.New("createFn is nil")
	}

	var createdTr *entities.Tree
	err := r.store.WithTx(ctx, func(s *store.Store) error {
		oldStore := r.store
		defer func() {
			r.store = oldStore
		}()
		r.store = s

		entity := defaultTree()
		created, err := createFn(&entity)
		if err != nil {
			return err
		}
		if !created {
			return nil
		}

		if err := r.validateTreeEntity(&entity); err != nil {
			return err
		}

		id, err := r.createEntity(ctx, &entity)
		if err != nil {
			return err
		}

		createdTr, err = r.GetByID(ctx, id)
		if err != nil {
			return err
		}

		if handler != nil {
			if err := handler(ctx, createdTr, id); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdTr, nil
}

func (r *TreeRepository) createEntity(ctx context.Context, entity *entities.Tree) (int32, error) {
	var treeClusterID *int32
	if entity.TreeCluster != nil {
		treeClusterID = &entity.TreeCluster.ID
	}

	var sensorID *string
	if entity.Sensor != nil {
		sensorID = &entity.Sensor.ID
		if err := r.store.UnlinkSensorIDFromTrees(ctx, sensorID); err != nil {
			return -1, err
		}
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

func (r *TreeRepository) validateTreeEntity(tree *entities.Tree) error {
	if tree == nil {
		return errors.New("tree is nil")
	}
	if tree.Latitude < -90 || tree.Latitude > 90 {
		return storage.ErrInvalidLatitude
	}
	if tree.Longitude < -180 || tree.Longitude > 180 {
		return storage.ErrInvalidLongitude
	}
	return nil
}

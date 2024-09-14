package tree

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

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

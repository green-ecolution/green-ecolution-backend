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

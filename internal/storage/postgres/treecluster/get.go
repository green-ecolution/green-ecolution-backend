package treecluster

import (
	"context"
	"errors"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/jackc/pgx/v5"
)

func (r *TreeClusterRepository) GetAll(ctx context.Context) ([]*entities.TreeCluster, error) {
	rows, err := r.store.GetAllTreeClusters(ctx)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	data := r.mapper.FromSqlList(rows)
	for _, f := range data {
		f.Region, err = r.GetRegionByTreeClusterID(ctx, f.ID)
		if err != nil {
			return nil, err
		}

		f.Trees, err = r.GetLinkedTreesByTreeClusterID(ctx, f.ID)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (r *TreeClusterRepository) GetByID(ctx context.Context, id int32) (*entities.TreeCluster, error) {
	row, err := r.store.GetTreeClusterByID(ctx, id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	data := r.mapper.FromSql(row)
	data.Region, err = r.GetRegionByTreeClusterID(ctx, id)
	if err != nil {
		return nil, err
	}

	data.Trees, err = r.GetLinkedTreesByTreeClusterID(ctx, id)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *TreeClusterRepository) GetSensorByTreeClusterID(ctx context.Context, id int32) (*entities.Sensor, error) {
	row, err := r.store.GetSensorByTreeClusterID(ctx, id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	return r.sensorMapper.FromSql(row), nil
}

func (r *TreeClusterRepository) GetRegionByTreeClusterID(ctx context.Context, id int32) (*entities.Region, error) {
	row, err := r.store.GetRegionByTreeClusterID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrRegionNotFound
		}
		return nil, r.store.HandleError(err)
	}

	return r.regionMapper.FromSql(row), nil
}

func (r *TreeClusterRepository) GetLinkedTreesByTreeClusterID(ctx context.Context, id int32) ([]*entities.Tree, error) {
	rows, err := r.store.GetLinkedTreesByTreeClusterID(ctx, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrRegionNotFound
		}
		return nil, r.store.HandleError(err)
	}

	return r.treeMapper.FromSqlList(rows), nil
}

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
	for _, tc := range data {
		if err := r.mapFields(ctx, tc); err != nil {
			return nil, r.store.HandleError(err)
		}
	}

	return data, nil
}

func (r *TreeClusterRepository) GetByID(ctx context.Context, id int32) (*entities.TreeCluster, error) {
	row, err := r.store.GetTreeClusterByID(ctx, id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	tc := r.mapper.FromSql(row)
	if err := r.mapFields(ctx, tc); err != nil {
		return nil, r.store.HandleError(err)
	}

	return tc, nil
}

func (r *TreeClusterRepository) GetByIDs(ctx context.Context, ids []int32) ([]*entities.TreeCluster, error) {
	rows, err := r.store.GetTreesClustersByIDs(ctx, ids)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	tc := r.mapper.FromSqlList(rows)
	for _, cluster := range tc {
		if err := r.mapFields(ctx, cluster); err != nil {
			return nil, r.store.HandleError(err)
		}
	}

	return tc, nil
}

func (r *TreeClusterRepository) GetRegionByTreeClusterID(ctx context.Context, id int32) (*entities.Region, error) {
	if err := r.tcIDExists(ctx, id); err != nil {
		return nil, err
	}

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
	if err := r.tcIDExists(ctx, id); err != nil {
		return nil, err
	}

	rows, err := r.store.GetLinkedTreesByTreeClusterID(ctx, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*entities.Tree{}, nil
		}
		return nil, r.store.HandleError(err)
	}

	return r.treeMapper.FromSqlList(rows), nil
}

func (r *TreeClusterRepository) tcIDExists(ctx context.Context, id int32) error {
	_, err := r.store.GetTreeClusterByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.ErrTreeClusterNotFound
		}
		return err
	}

	return nil
}

func (r *TreeClusterRepository) mapFields(ctx context.Context, tc *entities.TreeCluster) error {
	if err := r.mapRegion(ctx, tc); err != nil {
		return r.store.HandleError(err)
	}

	if err := r.mapTrees(ctx, tc); err != nil {
		return r.store.HandleError(err)
	}

	return nil
}

func (r *TreeClusterRepository) mapRegion(ctx context.Context, tc *entities.TreeCluster) error {
	region, err := r.GetRegionByTreeClusterID(ctx, tc.ID)
	if err != nil {
		// If region is not found, we can still return the tree cluster
		if !errors.Is(err, storage.ErrRegionNotFound) {
			return err
		}
	}
	tc.Region = region

	return nil
}

func (r *TreeClusterRepository) mapTrees(ctx context.Context, tc *entities.TreeCluster) error {
	trees, err := r.GetLinkedTreesByTreeClusterID(ctx, tc.ID)
	if err != nil {
		return err
	}
	tc.Trees = trees

	return nil
}

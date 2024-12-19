package treecluster

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
)

func (r *TreeClusterRepository) GetAll(ctx context.Context) ([]*entities.TreeCluster, error) {
	rows, err := r.store.GetAllTreeClusters(ctx)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	data := r.mapper.FromSqlList(rows)
	for _, tc := range data {
		if err := r.store.MapClusterFields(ctx, tc); err != nil {
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
	if err := r.store.MapClusterFields(ctx, tc); err != nil {
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
		if err := r.store.MapClusterFields(ctx, cluster); err != nil {
			return nil, r.store.HandleError(err)
		}
	}

	return tc, nil
}

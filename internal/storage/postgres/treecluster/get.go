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

	return r.mapper.FromSqlList(rows), nil
}

func (r *TreeClusterRepository) GetByID(ctx context.Context, id int32) (*entities.TreeCluster, error) {
	row, err := r.store.GetTreeClusterByID(ctx, id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	return r.mapper.FromSql(row), nil
}

func (r *TreeClusterRepository) GetSensorByTreeClusterID(ctx context.Context, id int32) (*entities.Sensor, error) {
	row, err := r.store.GetSensorByTreeClusterID(ctx, id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	return r.sensorMapper.FromSql(row), nil
}

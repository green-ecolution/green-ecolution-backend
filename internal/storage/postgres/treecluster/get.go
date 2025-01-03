package treecluster

import (
	"context"
	"errors"
	"fmt"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/twpayne/go-geos"
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

func (r *TreeClusterRepository) GetCenterPoint(ctx context.Context, tcID int32) (lat, long float64, err error) {
	geoStr, err := r.store.CalculateTreesCentroid(ctx, &tcID)
	if err != nil {
		return 0, 0, err
	}

	// Parse geoStr to get latitude and longitude
	g, err := geos.NewGeomFromWKT(geoStr)
	if err != nil {
		return 0, 0, err
	}

	if g.IsEmpty() {
		return 0, 0, errors.New("empty geometry")
	}

	return g.X(), g.Y(), nil
}

func (r *TreeClusterRepository) GetAllLatestSensorDataByClusterID(ctx context.Context, tcID int32) ([]*entities.SensorData, error) {
	rows, err := r.store.GetAllLatestSensorDataByTreeClusterID(ctx, tcID)
	if err != nil {
		return nil, r.store.HandleError(err)
	}
	domainData, err := r.sensorMapper.FromSqlSensorDataList(rows)
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("failed to map sensor data"))
	}

	return domainData, nil
}

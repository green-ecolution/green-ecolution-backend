package treecluster

import (
	"context"
	"errors"
	"fmt"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/twpayne/go-geos"
)

func (r *TreeClusterRepository) GetAll(ctx context.Context) ([]*entities.TreeCluster, int64, error) {
	log := logger.GetLogger(ctx)
	rows, err := r.store.GetAllTreeClusters(ctx, &sqlc.GetAllTreeClustersParams{
		Limit: limit,
		Offset: (page - 1) * limit,
	})

	if err != nil {
		log.Debug("failed to get tree clusters in db")
		return nil, 0, r.store.MapError(err, sqlc.TreeCluster{})
	}

	data := r.mapper.FromSqlList(rows)
	for _, tc := range data {
		if err := r.store.MapClusterFields(ctx, tc); err != nil {
			return nil, 0, r.store.MapError(err, sqlc.TreeCluster{})
		}
	}

	return data, totalCount, nil
}

func (r *TreeClusterRepository) GetByID(ctx context.Context, id int32) (*entities.TreeCluster, error) {
	log := logger.GetLogger(ctx)
	row, err := r.store.GetTreeClusterByID(ctx, id)
	if err != nil {
		log.Debug("failed to get tree cluster by id in db", "error", err, "cluster_id", id)
		return nil, r.store.MapError(err, sqlc.TreeCluster{})
	}

	tc := r.mapper.FromSql(row)
	if err := r.store.MapClusterFields(ctx, tc); err != nil {
		return nil, r.store.MapError(err, sqlc.TreeCluster{})
	}

	return tc, nil
}

func (r *TreeClusterRepository) GetByIDs(ctx context.Context, ids []int32) ([]*entities.TreeCluster, error) {
	log := logger.GetLogger(ctx)
	rows, err := r.store.GetTreesClustersByIDs(ctx, ids)
	if err != nil {
		log.Debug("failed to get tree cluster by multiple ids", "error", err, "cluster_ids", ids)
		return nil, r.store.MapError(err, sqlc.TreeCluster{})
	}

	tc := r.mapper.FromSqlList(rows)
	for _, cluster := range tc {
		if err := r.store.MapClusterFields(ctx, cluster); err != nil {
			return nil, r.store.MapError(err, sqlc.TreeCluster{})
		}
	}

	return tc, nil
}

func (r *TreeClusterRepository) GetCenterPoint(ctx context.Context, tcID int32) (lat, long float64, err error) {
	log := logger.GetLogger(ctx)
	geoStr, err := r.store.CalculateTreesCentroid(ctx, &tcID)
	if err != nil {
		log.Warn("failed to calculate center point of given cluster", "error", err, "cluster_id", tcID)
		return 0, 0, err
	}

	// Parse geoStr to get latitude and longitude
	g, err := geos.NewGeomFromWKT(geoStr)
	if err != nil {
		log.Debug("failed to parse calculated geo string", "error", err, "geo_string", geoStr)
		return 0, 0, err
	}

	if g.IsEmpty() {
		return 0, 0, errors.New("empty geometry")
	}

	return g.X(), g.Y(), nil
}

func (r *TreeClusterRepository) GetAllLatestSensorDataByClusterID(ctx context.Context, tcID int32) ([]*entities.SensorData, error) {
	log := logger.GetLogger(ctx)
	rows, err := r.store.GetAllLatestSensorDataByTreeClusterID(ctx, tcID)
	if err != nil {
		log.Debug("failed to get all latest sensor data by given cluster id in db", "error", err, "cluster_id", tcID)
		return nil, r.store.MapError(err, sqlc.TreeCluster{})
	}
	domainData, err := r.sensorMapper.FromSqlSensorDataList(rows)
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("failed to map sensor data"))
	}

	return domainData, nil
}

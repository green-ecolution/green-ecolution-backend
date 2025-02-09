package treecluster

import (
	"context"
	"errors"
	"fmt"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils/pagination"
	"github.com/twpayne/go-geos"
)

func (r *TreeClusterRepository) GetAll(ctx context.Context, provider string) ([]*entities.TreeCluster, int64, error) {
	log := logger.GetLogger(ctx)
	page, limit, err := pagination.GetValues(ctx)
	if err != nil {
		return nil, 0, r.store.MapError(err, sqlc.TreeCluster{})
	}

	totalCount, err := r.store.GetAllTreeClustersCount(ctx, provider)
	if err != nil {
		log.Debug("failed to get total tree cluster count in db", "error", err)
		return nil, 0, r.store.MapError(err, sqlc.TreeCluster{})
	}

	if totalCount == 0 {
		return []*entities.TreeCluster{}, 0, nil
	}

	if limit == -1 {
		limit = int32(totalCount)
		page = 1
	}

	rows, err := r.store.GetAllTreeClusters(ctx, &sqlc.GetAllTreeClustersParams{
		Column1: provider,
		Limit:   limit,
		Offset:  (page - 1) * limit,
	})

	if err != nil {
		log.Debug("failed to get tree clusters in db")
		return nil, 0, r.store.MapError(err, sqlc.TreeCluster{})
	}

	data, err := r.mapper.FromSqlList(rows)
	if err != nil {
		log.Debug("failed to convert entity", "error", err)
		return nil, 0, err
	}

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

	tc, err := r.mapper.FromSql(row)
	if err != nil {
		log.Debug("failed to convert entity", "error", err)
		return nil, err
	}

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

	tc, err := r.mapper.FromSqlList(rows)
	if err != nil {
		log.Debug("failed to convert entity", "error", err)
		return nil, err
	}

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

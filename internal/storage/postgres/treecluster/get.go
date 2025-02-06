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

func (r *TreeClusterRepository) GetAll(ctx context.Context, filter entities.TreeClusterFilter) ([]*entities.TreeCluster, error) {
	log := logger.GetLogger(ctx)

	rows, err := r.store.GetAllTreeClusters(ctx, &sqlc.GetAllTreeClustersParams{
		Column1: sqlc.WateringStatus(filter.WateringStatus),
		Column2: filter.Region,
		Limit:   filter.Limit,
		Offset:  filter.Offset,
	})

	if err != nil {
		log.Debug("failed to get tree clusters in db")
		return nil, r.store.MapError(err, sqlc.TreeCluster{})
	}

	data, err := r.mapper.FromSqlList(rows)
	if err != nil {
		log.Debug("failed to convert entity", "error", err)
		return nil, err
	}

	for _, tc := range data {
		if err := r.store.MapClusterFields(ctx, tc); err != nil {
			return nil, r.store.MapError(err, sqlc.TreeCluster{})
		}
	}

	return data, nil
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

func (r *TreeClusterRepository) GetTreeClustersCount(ctx context.Context) (int64, error) {
	log := logger.GetLogger(ctx)
	count, err := r.store.GetAllTreeClustersCount(ctx)
	if err != nil {
		log.Debug("Failed to get total tree cluster count", "error", err)
		return 0, errors.New("failed to get the total tree cluster count")
	}
	return count, nil
}

func (r *TreeClusterRepository) GetTreeClustersCountByStatus(ctx context.Context, status entities.WateringStatus) (int64, error) {
	log := logger.GetLogger(ctx)
	count, err := r.store.GetTreeClustersCountByStatus(ctx, sqlc.WateringStatus(status))
	if err != nil {
		log.Debug("Failed to get count by status", "error", err)
		return 0, errors.New("failed to get count by status")
	}
	return count, nil
}

func (r *TreeClusterRepository) GetTreeClustersCountByRegion(ctx context.Context, region string) (int64, error) {
	log := logger.GetLogger(ctx)
	count, err := r.store.GetTreeClustersCountByRegion(ctx, region)
	if err != nil {
		log.Debug("Failed to get count by region", "error", err)
		return 0, errors.New("failed to get count by region")
	}
	return count, nil
}

func (r *TreeClusterRepository) GetTreeClustersCountByStatusAndRegion(ctx context.Context, status entities.WateringStatus, region string) (int64, error) {
	log := logger.GetLogger(ctx)
	count, err := r.store.GetTreeClustersCountByStatusAndRegion(ctx, &sqlc.GetTreeClustersCountByStatusAndRegionParams{
		WateringStatus: sqlc.WateringStatus(status),
		Name:           region,
	})
	if err != nil {
		log.Debug("Failed to get count by status and region", "error", err)
		return 0, errors.New("failed to get count by status and region")
	}
	return count, nil
}

package treecluster

import (
	"context"
	"log/slog"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
)

type GeoClusterLocator struct {
	treeRepo   storage.TreeRepository
	regionRepo storage.RegionRepository
}

func NewGeoLocation(treeRepo storage.TreeRepository, regionRepo storage.RegionRepository) *GeoClusterLocator {
	return &GeoClusterLocator{
		treeRepo:   treeRepo,
		regionRepo: regionRepo,
	}
}

func (s *GeoClusterLocator) UpdateCluster(ctx context.Context, cluster *entities.TreeCluster) error {
	if cluster == nil {
		return nil
	}

	slog.Debug("Updating cluster location", "clusterID", cluster.ID)
	if len(cluster.Trees) == 0 {
		s.removeClusterCoords(cluster)
		return nil
	}

	return s.setClusterCoords(ctx, cluster, cluster.Trees)
}

func (s *GeoClusterLocator) removeClusterCoords(tc *entities.TreeCluster) {
	tc.Latitude = nil
	tc.Longitude = nil
	tc.Region = nil
}

func (s *GeoClusterLocator) setClusterCoords(ctx context.Context, tc *entities.TreeCluster, trees []*entities.Tree) error {
	treeIDs := utils.Map(trees, func(t *entities.Tree) int32 {
		return t.ID
	})

	lat, long, err := s.calculateCenterPoint(ctx, treeIDs)
	if err != nil {
		return err
	}

	region, err := s.getRegionByPoint(ctx, lat, long)
	if err != nil {
		return err
	}

	tc.Latitude = &lat
	tc.Longitude = &long
	tc.Region = region

	return nil
}

func (s *GeoClusterLocator) getRegionByPoint(ctx context.Context, lat, long float64) (*entities.Region, error) {
	region, err := s.regionRepo.GetByPoint(ctx, lat, long)
	if err != nil {
		return nil, err
	}

	return region, nil
}

func (s *GeoClusterLocator) calculateCenterPoint(ctx context.Context, treeIDs []int32) (lat, long float64, err error) {
	lat, long, err = s.treeRepo.GetCenterPoint(ctx, treeIDs)
	if err != nil {
		return 0, 0, err
	}

	return lat, long, nil
}

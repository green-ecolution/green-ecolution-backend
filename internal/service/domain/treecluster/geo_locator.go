package treecluster

import (
	"context"
	"log/slog"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/treecluster"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
)

type GeoClusterLocator struct {
	treeRepo    storage.TreeRepository
	clusterRepo storage.TreeClusterRepository
	regionRepo  storage.RegionRepository
}

func NewLocationUpdate(clusterRepo storage.TreeClusterRepository, treeRepo storage.TreeRepository, regionRepo storage.RegionRepository) *GeoClusterLocator {
	return &GeoClusterLocator{
		clusterRepo: clusterRepo,
		treeRepo:    treeRepo,
		regionRepo:  regionRepo,
	}
}

// UpdateCluster updates the center point of a cluster based on the center point of its trees
func (s *GeoClusterLocator) UpdateCluster(ctx context.Context, clusterID *int32) error {
	slog.Debug("Updating cluster location", "clusterID", clusterID)
	if clusterID == nil {
		return nil
	}
	cluster, err := s.clusterRepo.GetByID(ctx, *clusterID)
	if err != nil {
		return err
	}

	if len(cluster.Trees) == 0 {
		return s.removeClusterCoords(ctx, *clusterID)
	}

	return s.setClusterCoords(ctx, *clusterID, cluster.Trees)
}

func (s *GeoClusterLocator) removeClusterCoords(ctx context.Context, clusterID int32) error {
	_, err := s.clusterRepo.Update(ctx, clusterID,
		treecluster.WithLatitude(nil),
		treecluster.WithLongitude(nil),
		treecluster.WithRegion(nil),
	)
	return err
}

func (s *GeoClusterLocator) setClusterCoords(ctx context.Context, clusterID int32, trees []*entities.Tree) error {
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

	_, err = s.clusterRepo.Update(ctx, clusterID,
		treecluster.WithLatitude(&lat),
		treecluster.WithLongitude(&long),
		treecluster.WithRegion(region),
	)

	return err
}

func (s *GeoClusterLocator) getRegionByPoint(ctx context.Context, lat, long float64) (*entities.Region, error) {
	region, err := s.regionRepo.GetByPoint(ctx, lat, long)
	if err != nil {
		return nil, handleError(err)
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

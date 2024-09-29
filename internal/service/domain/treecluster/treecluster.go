package treecluster

import (
	"context"
	"errors"

	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/treecluster"
)

type TreeClusterService struct {
	treeClusterRepo storage.TreeClusterRepository
	treeRepo        storage.TreeRepository
	regionRepo      storage.RegionRepository
}

func NewTreeClusterService(treeClusterRepo storage.TreeClusterRepository, treeRepo storage.TreeRepository, regionRepo storage.RegionRepository) service.TreeClusterService {
	return &TreeClusterService{
		treeClusterRepo: treeClusterRepo,
		treeRepo:        treeRepo,
		regionRepo:      regionRepo,
	}
}

func (s *TreeClusterService) GetAll(ctx context.Context) ([]*domain.TreeCluster, error) {
	treeClusters, err := s.treeClusterRepo.GetAll(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	return treeClusters, nil
}

func (s *TreeClusterService) GetByID(ctx context.Context, id int32) (*domain.TreeCluster, error) {
	treeCluster, err := s.treeClusterRepo.GetByID(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}

	return treeCluster, nil
}

func (s *TreeClusterService) Create(ctx context.Context, tc *domain.TreeClusterCreate) (*domain.TreeCluster, error) {
  treeIDs := make([]int32, len(tc.TreeIDs))
  for i, treeID := range tc.TreeIDs {
    treeIDs[i] = *treeID
  }

	lat, long, err := s.treeRepo.GetCenterPoint(ctx, treeIDs)
	if err != nil {
		return nil, err
	}

	region, err := s.regionRepo.GetByPoint(ctx, lat, long)
	if err != nil {
		return nil, err
	}

	c, err := s.treeClusterRepo.Create(ctx,
		treecluster.WithName(tc.Name),
		treecluster.WithAddress(tc.Address),
		treecluster.WithDescription(tc.Description),
		treecluster.WithLatitude(lat),
		treecluster.WithLongitude(long),
		treecluster.WithRegion(region),
	)
	if err != nil {
		return nil, handleError(err)
	}

  if err = s.treeRepo.UpdateTreeClusterID(ctx, treeIDs, &c.ID); err != nil {
    return nil, handleError(err)
  }

	return c, nil
}

func (s *TreeClusterService) Delete(ctx context.Context, id int32) error {
	_, err := s.treeClusterRepo.GetByID(ctx, id)
	if err != nil {
		return handleError(err)
	}

	// TODO: Add a transaction to undo this change if an error occurs.
	err = s.treeRepo.UnlinkTreeClusterID(ctx, id)
	if err != nil {
		return handleError(err)
	}

	err = s.treeClusterRepo.Delete(ctx, id)
	if err != nil {
		return handleError(err)
	}

	return nil
}

func handleError(err error) error {
	if errors.Is(err, storage.ErrEntityNotFound) {
		return service.NewError(service.NotFound, err.Error())
	}

	return service.NewError(service.InternalError, err.Error())
}

func (s *TreeClusterService) Ready() bool {
	return s.treeClusterRepo != nil
}

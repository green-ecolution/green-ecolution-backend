package treecluster

import (
	"context"
	"errors"
	"log/slog"

	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/treecluster"
)

type TreeClusterService struct {
	treeClusterRepo storage.TreeClusterRepository
	treeRepo        storage.TreeRepository
	regionRepo      storage.RegionRepository
	locator         *GeoClusterLocator
}

func NewTreeClusterService(
	treeClusterRepo storage.TreeClusterRepository,
	treeRepo storage.TreeRepository,
	regionRepo storage.RegionRepository,
	locator *GeoClusterLocator,
) service.TreeClusterService {
	return &TreeClusterService{
		treeClusterRepo: treeClusterRepo,
		treeRepo:        treeRepo,
		regionRepo:      regionRepo,
		locator:         locator,
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
	trees, err := s.getTrees(ctx, tc.TreeIDs)
	if err != nil {
		return nil, handleError(err)
	}

	visitedClusters := make(map[int32]bool)
	for _, tree := range trees {
		if tree.TreeCluster != nil && tree.TreeCluster.ID != 0 {
			if _, ok := visitedClusters[tree.TreeCluster.ID]; ok {
				slog.Debug("Tree already visited", "treeID", tree.ID)
				continue
			}

			slog.Debug("Updating cluster", "clusterID", tree.TreeCluster.ID)
			if err = s.locator.UpdateCluster(ctx, &tree.TreeCluster.ID); err != nil {
				return nil, handleError(err)
			}
			visitedClusters[tree.TreeCluster.ID] = true
		}
	}

	c, err := s.treeClusterRepo.Create(ctx,
		treecluster.WithName(tc.Name),
		treecluster.WithAddress(tc.Address),
		treecluster.WithDescription(tc.Description),
		treecluster.WithTrees(trees),
	)
	if err != nil {
		return nil, handleError(err)
	}

	if err = s.locator.UpdateCluster(ctx, &c.ID); err != nil {
		return nil, handleError(err)
	}

	return c, nil
}

func (s *TreeClusterService) Update(ctx context.Context, id int32, tc *domain.TreeClusterUpdate) (*domain.TreeCluster, error) {
	trees, err := s.getTrees(ctx, tc.TreeIDs)
	if err != nil {
		return nil, handleError(err)
	}

	c, err := s.treeClusterRepo.Update(ctx, id,
		treecluster.WithTrees(trees),
		treecluster.WithName(tc.Name),
		treecluster.WithAddress(tc.Address),
		treecluster.WithDescription(tc.Description),
		treecluster.WithSoilCondition(tc.SoilCondition),
	)

	if err != nil {
		return nil, handleError(err)
	}

	if err = s.locator.UpdateCluster(ctx, &id); err != nil {
		return nil, handleError(err)
	}

	return c, nil
}

func (s *TreeClusterService) Delete(ctx context.Context, id int32) error {
	_, err := s.treeClusterRepo.GetByID(ctx, id)
	if err != nil {
		return handleError(err)
	}

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

func (s *TreeClusterService) Ready() bool {
	return s.treeClusterRepo != nil
}

func handleError(err error) error {
	// TODO: Rollback the transaction if an error occurs.
	if errors.Is(err, storage.ErrEntityNotFound) {
		return service.NewError(service.NotFound, err.Error())
	}

	return service.NewError(service.InternalError, err.Error())
}

func (s *TreeClusterService) getTrees(ctx context.Context, ids []*int32) ([]*domain.Tree, error) {
	treeIDs := make([]int32, len(ids))
	for i, id := range ids {
		treeIDs[i] = *id
	}

	var err error
	trees, err := s.treeRepo.GetTreesByIDs(ctx, treeIDs)
	if err != nil {
		return nil, err
	}

	return trees, nil
}

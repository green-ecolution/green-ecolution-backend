package treecluster

import (
	"context"
	"log/slog"

	"github.com/pkg/errors"

	"github.com/go-playground/validator/v10"
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
)

type TreeClusterService struct {
	treeClusterRepo storage.TreeClusterRepository
	treeRepo        storage.TreeRepository
	regionRepo      storage.RegionRepository
	locator         service.GeoClusterLocator
	validator       *validator.Validate
}

func NewTreeClusterService(
	treeClusterRepo storage.TreeClusterRepository,
	treeRepo storage.TreeRepository,
	regionRepo storage.RegionRepository,
	locator service.GeoClusterLocator,
) service.TreeClusterService {
	return &TreeClusterService{
		treeClusterRepo: treeClusterRepo,
		treeRepo:        treeRepo,
		regionRepo:      regionRepo,
		locator:         locator,
		validator:       validator.New(),
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

func (s *TreeClusterService) Create(ctx context.Context, createTc *domain.TreeClusterCreate) (*domain.TreeCluster, error) {
	if err := s.validator.Struct(createTc); err != nil {
		return nil, service.NewError(service.BadRequest, errors.Wrap(err, "validation error").Error())
	}

	trees, err := s.getTrees(ctx, createTc.TreeIDs)
	if err != nil {
		return nil, handleError(err)
	}

	c, err := s.treeClusterRepo.Create(ctx, func(tc *domain.TreeCluster) (bool, error) {
		if err = s.handlePrevTreeLocation(ctx, trees); err != nil {
			return false, handleError(err)
		}

		tc.Name = createTc.Name
		tc.Address = createTc.Address
		tc.Description = createTc.Description
		tc.Trees = trees

		if err = s.locator.UpdateCluster(ctx, tc); err != nil {
			return false, handleError(err)
		}

		return true, nil
	})

	if err != nil {
		return nil, handleError(err)
	}

	return c, nil
}

func (s *TreeClusterService) Update(ctx context.Context, id int32, tcUpdate *domain.TreeClusterUpdate) (*domain.TreeCluster, error) {
	if err := s.validator.Struct(tcUpdate); err != nil {
		return nil, service.NewError(service.BadRequest, errors.Wrap(err, "validation error").Error())
	}

	trees, err := s.getTrees(ctx, tcUpdate.TreeIDs)
	if err != nil {
		return nil, handleError(err)
	}

	err = s.treeClusterRepo.Update(ctx, id, func(tc *domain.TreeCluster) (bool, error) {
		tc.Trees = trees
		tc.Name = tcUpdate.Name
		tc.Address = tcUpdate.Address
		tc.Description = tcUpdate.Description
		tc.SoilCondition = tcUpdate.SoilCondition

		if err := s.locator.UpdateCluster(ctx, tc); err != nil {
			return false, err
		}

		return true, nil
	})

	if err != nil {
		return nil, handleError(err)
	}

	return s.GetByID(ctx, id)
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

// handlePrevTreeLocation updates the locations of clusters associated with the provided trees.
//
// This method iterates over a list of trees and processes the clusters they belong to.
// For each cluster, the cluster's location is updated by recalculating its coordinates and region using the `GeoClusterLocator`.
// Clusters are only updated once, even if multiple trees belong to the same cluster.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation, and timeouts.
//   - trees: A slice of Tree entities to process. Each tree may reference a cluster.
//
// Returns:
//
//	An error if any cluster update fails, otherwise nil.
//
// Notes:
//   - Clusters with an ID of 0 are ignored.
//   - Updates are performed via a callback mechanism in the `treeClusterRepo` to ensure thread safety or transactional consistency.
func (s *TreeClusterService) handlePrevTreeLocation(ctx context.Context, trees []*domain.Tree) error {
	visitedClusters := make(map[int32]bool)
	for _, tree := range trees {
		if tree.TreeCluster != nil && tree.TreeCluster.ID != 0 {
			continue
		}

		if _, ok := visitedClusters[tree.TreeCluster.ID]; ok {
			slog.Debug("Tree already visited", "treeID", tree.ID)
			continue
		}

		slog.Debug("Updating cluster", "clusterID", tree.TreeCluster.ID)

		err := s.treeClusterRepo.Update(ctx, tree.TreeCluster.ID, func(tc *domain.TreeCluster) (bool, error) {
			if err := s.locator.UpdateCluster(ctx, tc); err != nil {
				return false, err
			}
			return true, nil
		})
		if err != nil {
			return err
		}

		visitedClusters[tree.TreeCluster.ID] = true
	}

	return nil
}

func handleError(err error) error {
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

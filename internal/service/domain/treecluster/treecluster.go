package treecluster

import (
	"context"
	"log/slog"

	"github.com/pkg/errors"

	"github.com/go-playground/validator/v10"
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/worker"
)

type TreeClusterService struct {
	treeClusterRepo storage.TreeClusterRepository
	treeRepo        storage.TreeRepository
	regionRepo      storage.RegionRepository
	validator       *validator.Validate
	eventManager    *worker.EventManager
}

func NewTreeClusterService(
	treeClusterRepo storage.TreeClusterRepository,
	treeRepo storage.TreeRepository,
	regionRepo storage.RegionRepository,
	eventManager *worker.EventManager,
) service.TreeClusterService {
	return &TreeClusterService{
		treeClusterRepo: treeClusterRepo,
		treeRepo:        treeRepo,
		regionRepo:      regionRepo,
		validator:       validator.New(),
		eventManager:    eventManager,
	}
}

type TreeClusterOption func(*entities.TreeCluster)

func WithName(name string) TreeClusterOption {
    return func(tc *entities.TreeCluster) {
        tc.Name = name
    }
}

func WithAddress(address string) TreeClusterOption {
    return func(tc *entities.TreeCluster) {
        tc.Address = address
    }
}

func WithDescription(description string) TreeClusterOption {
    return func(tc *entities.TreeCluster) {
        tc.Description = description
    }
}

func WithMoistureLevel(moistureLevel float64) TreeClusterOption {
    return func(tc *entities.TreeCluster) {
        tc.MoistureLevel = moistureLevel
    }
}

func WithWateringStatus(wateringStatus entities.WateringStatus) TreeClusterOption {
    return func(tc *entities.TreeCluster) {
        tc.WateringStatus = wateringStatus
    }
}

func WithSoilCondition(soilCondition entities.TreeSoilCondition) TreeClusterOption {
    return func(tc *entities.TreeCluster) {
        tc.SoilCondition = soilCondition
    }
}

func WithArchived(archived bool) TreeClusterOption {
    return func(tc *entities.TreeCluster) {
        tc.Archived = archived
    }
}

func WithLastWatered(lastWatered time.Time) TreeClusterOption {
    return func(tc *entities.TreeCluster) {
        tc.LastWatered = lastWatered
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

func (s *TreeClusterService) getUpdatedLatLong(ctx context.Context, tc *domain.TreeCluster) (lat, long *float64, region *domain.Region, err error) {
	if len(tc.Trees) == 0 {
		return nil, nil, nil, nil
	}

	latitude, longitude, err := s.treeClusterRepo.GetCenterPoint(ctx, tc.ID)
	if err != nil {
		slog.Error("failed to get center point of tree cluster", "error", err, "tree_cluster", tc)
		return nil, nil, nil, err
	}

	region, err = s.regionRepo.GetByPoint(ctx, latitude, longitude)
	if err != nil {
		slog.Error("can't find region by lat and long", "error", err, "latitude", latitude, "longitude", longitude, "tree_cluster", tc)
		return &latitude, &longitude, nil, nil
	}

	return &latitude, &longitude, region, nil
}

func (s *TreeClusterService) publishUpdateEvent(ctx context.Context, prevTc *domain.TreeCluster) error {
	slog.Debug("publish new event", "event", domain.EventTypeUpdateTreeCluster, "service", "TreeClusterService")
	updatedTc, err := s.GetByID(ctx, prevTc.ID)
	if err != nil {
		return err
	}
	event := domain.NewEventUpdateTreeCluster(prevTc, updatedTc)
	err = s.eventManager.Publish(ctx, event)
	if err != nil {
		slog.Error("error while sending event after updating tree cluster", "err", err)
	}

	return nil
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
			slog.Error("failed to update prev tree location", "error", err, "trees", trees, "tree_cluster", tc)
			return false, handleError(err)
		}

		tc.Trees = trees
		tc.Name = createTc.Name
		tc.Address = createTc.Address
		tc.Description = createTc.Description
		tc.SoilCondition = createTc.SoilCondition

		return true, nil
	})

	if err != nil {
		return nil, handleError(err)
	}

	if err := s.updateTreeClusterPosition(ctx, c.ID); err != nil {
		slog.Error("error while update the cluster locations", "error", err, "cluster_id", c.ID)
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

	prevTc, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}

	err = s.treeClusterRepo.Update(ctx, id, func(tc *domain.TreeCluster) (bool, error) {
		tc.Trees = trees
		tc.Name = tcUpdate.Name
		tc.Address = tcUpdate.Address
		tc.Description = tcUpdate.Description
		tc.SoilCondition = tcUpdate.SoilCondition

		return true, nil
	})

	if err != nil {
		return nil, handleError(err)
	}

	if err := s.updateTreeClusterPosition(ctx, id); err != nil {
		slog.Error("error while update the cluster locations", "error", err, "cluster_id", id)
		return nil, handleError(err)
	}

	if err := s.publishUpdateEvent(ctx, prevTc); err != nil {
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

// Update the tree cluster only after the trees have been updated to the database,
// otherwise the center point of the tree cluster cannot be set
func (s *TreeClusterService) updateTreeClusterPosition(ctx context.Context, id int32) error {
	err := s.treeClusterRepo.Update(ctx, id, func(tc *domain.TreeCluster) (bool, error) {
		lat, long, region, err := s.getUpdatedLatLong(ctx, tc)
		if err != nil {
			return false, nil
		}

		tc.Latitude = lat
		tc.Longitude = long
		tc.Region = region

		return true, nil
	})

	return err
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
		if tree.TreeCluster == nil || tree.TreeCluster.ID == 0 {
			continue
		}

		if _, ok := visitedClusters[tree.TreeCluster.ID]; ok {
			slog.Debug("Tree already visited", "treeID", tree.ID)
			continue
		}

		slog.Debug("Updating cluster", "clusterID", tree.TreeCluster.ID)

		updateFn := func(tc *domain.TreeCluster) (bool, error) {
			lat, long, region, err := s.getUpdatedLatLong(ctx, tc)
			if err != nil {
				return false, err
			}

			tc.Latitude = lat
			tc.Longitude = long
			tc.Region = region

			return true, nil
		}

		if err := s.treeClusterRepo.Update(ctx, tree.TreeCluster.ID, updateFn); err != nil {
			return err
		}

		if err := s.publishUpdateEvent(ctx, tree.TreeCluster); err != nil {
			return err
		}

		visitedClusters[tree.TreeCluster.ID] = true
	}

	return nil
}

func handleError(err error) error {
	if errors.Is(err, storage.ErrEntityNotFound) {
		return service.NewError(service.NotFound, storage.ErrTreeClusterNotFound.Error())
	}

	return service.NewError(service.InternalError, err.Error())
}

func (s *TreeClusterService) getTrees(ctx context.Context, ids []*int32) ([]*domain.Tree, error) {
	treeIDs := make([]int32, len(ids))
	for i, id := range ids {
		treeIDs[i] = *id
	}

	return s.treeRepo.GetTreesByIDs(ctx, treeIDs)
}

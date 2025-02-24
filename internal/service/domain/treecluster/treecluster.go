package treecluster

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-playground/validator/v10"
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
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

func (s *TreeClusterService) GetAll(ctx context.Context, provider string) ([]*domain.TreeCluster, int64, error) {
	log := logger.GetLogger(ctx)
	treeClusters, totalCount, err := s.treeClusterRepo.GetAll(ctx, provider)
	if err != nil {
		log.Debug("failed to fetch tree clsuters", "error", err)
		return nil, 0, service.MapError(ctx, err, service.ErrorLogEntityNotFound)
	}
	return treeClusters, totalCount, nil
}

func (s *TreeClusterService) GetByID(ctx context.Context, id int32) (*domain.TreeCluster, error) {
	log := logger.GetLogger(ctx)
	treeCluster, err := s.treeClusterRepo.GetByID(ctx, id)
	if err != nil {
		log.Debug("failed to fetch tree cluster by id", "error", err, "cluster_id", id)
		return nil, service.MapError(ctx, err, service.ErrorLogEntityNotFound)
	}

	return treeCluster, nil
}

func (s *TreeClusterService) getUpdatedLatLong(ctx context.Context, tc *domain.TreeCluster) (lat, long *float64, region *domain.Region, err error) {
	log := logger.GetLogger(ctx)
	if len(tc.Trees) == 0 {
		return nil, nil, nil, nil
	}

	latitude, longitude, err := s.treeClusterRepo.GetCenterPoint(ctx, tc.ID)
	if err != nil {
		log.Error("failed to get center point of tree cluster", "error", err, "tree_cluster", tc)
		return nil, nil, nil, err
	}

	region, err = s.regionRepo.GetByPoint(ctx, latitude, longitude)
	if err != nil {
		log.Error("can't find region by lat and long", "error", err, "latitude", latitude, "longitude", longitude, "tree_cluster", tc)
		return &latitude, &longitude, nil, nil
	}

	return &latitude, &longitude, region, nil
}

func (s *TreeClusterService) publishUpdateEvent(ctx context.Context, prevTc *domain.TreeCluster) error {
	log := logger.GetLogger(ctx)
	log.Debug("publish new event", "event", domain.EventTypeUpdateTreeCluster, "service", "TreeClusterService")
	updatedTc, err := s.GetByID(ctx, prevTc.ID)
	if err != nil {
		return err
	}

	event := domain.NewEventUpdateTreeCluster(prevTc, updatedTc)
	err = s.eventManager.Publish(ctx, event)
	if err != nil {
		log.Error("error while sending event after updating tree cluster", "err", err)
	}

	return nil
}

func (s *TreeClusterService) Create(ctx context.Context, createTc *domain.TreeClusterCreate) (*domain.TreeCluster, error) {
	log := logger.GetLogger(ctx)
	if err := s.validator.Struct(createTc); err != nil {
		log.Debug("failed to validate struct in create tree cluster", "error", err, "raw_cluster", fmt.Sprintf("%+v", createTc))
		return nil, service.MapError(ctx, errors.Join(err, service.ErrValidation), service.ErrorLogValidation)
	}

	trees, err := s.getTrees(ctx, createTc.TreeIDs)
	if err != nil {
		log.Debug("failed to get trees inside the tree cluster", "error", err, "tree_ids", createTc.TreeIDs)
		return nil, service.MapError(ctx, err, service.ErrorLogEntityNotFound)
	}

	c, err := s.treeClusterRepo.Create(ctx, func(tc *domain.TreeCluster, repo storage.TreeClusterRepository) (bool, error) {

		if err = s.handlePrevTreeLocation(ctx, trees, repo.Update); err != nil {
			log.Debug("failed to update prev tree location", "error", err, "trees", trees, "tree_cluster", tc)
			return false, service.MapError(ctx, err, service.ErrorLogAll)
		}

		tc.Trees = trees
		tc.Name = createTc.Name
		tc.Address = createTc.Address
		tc.Description = createTc.Description
		tc.SoilCondition = createTc.SoilCondition
		tc.Provider = createTc.Provider
		tc.AdditionalInfo = createTc.AdditionalInfo

		log.Debug("creating tree cluster with following attributes",
			"tree_ids", createTc.TreeIDs,
			"name", createTc.Name,
			"address", createTc.Address,
			"description", createTc.Description,
			"soil_condition", createTc.SoilCondition,
		)

		return true, nil
	})

	if err != nil {
		log.Debug("failed to create tree cluster", "error", err)
		return nil, service.MapError(ctx, err, service.ErrorLogAll)
	}

	if err := s.updateTreeClusterPosition(ctx, c.ID); err != nil {
		log.Debug("error while update the cluster locations", "error", err, "cluster_id", c.ID)
		return nil, service.MapError(ctx, err, service.ErrorLogAll)
	}

	log.Info("tree cluster created successfully", "cluster_id", c.ID)
	return c, nil
}

func (s *TreeClusterService) Update(ctx context.Context, id int32, tcUpdate *domain.TreeClusterUpdate) (*domain.TreeCluster, error) {
	log := logger.GetLogger(ctx)
	if err := s.validator.Struct(tcUpdate); err != nil {
		log.Debug("failed to validate struct from update tree cluster request", "error", err, "raw_cluster", fmt.Sprintf("%+v", tcUpdate))
		return nil, service.MapError(ctx, errors.Join(err, service.ErrValidation), service.ErrorLogValidation)
	}

	trees, err := s.getTrees(ctx, tcUpdate.TreeIDs)
	if err != nil {
		log.Debug("failed to get trees inside the tree cluster", "error", err, "tree_ids", tcUpdate.TreeIDs)
		return nil, service.MapError(ctx, err, service.ErrorLogEntityNotFound)
	}

	prevTc, err := s.GetByID(ctx, id)
	if err != nil {
		log.Debug("failed to get exiting tree cluster", "error", err, "cluster_id", id)
		return nil, service.MapError(ctx, err, service.ErrorLogEntityNotFound)
	}

	err = s.treeClusterRepo.Update(ctx, id, func(tc *domain.TreeCluster) (bool, error) {
		tc.Trees = trees
		tc.Name = tcUpdate.Name
		tc.Address = tcUpdate.Address
		tc.Description = tcUpdate.Description
		tc.SoilCondition = tcUpdate.SoilCondition
		tc.Provider = tcUpdate.Provider
		tc.AdditionalInfo = tcUpdate.AdditionalInfo

		log.Debug("updating tree cluster with following attributes",
			"cluster_id", id,
			"name", tcUpdate.Name,
			"address", tcUpdate.Address,
			"description", tcUpdate.Description,
			"soil_condition", tcUpdate.SoilCondition,
			"provider", tcUpdate.Provider,
			"additional_info", tcUpdate.AdditionalInfo,
		)

		return true, nil
	})

	if err != nil {
		log.Debug("failed to update tree cluster", "error", err, "cluster_id", id)
		return nil, service.MapError(ctx, err, service.ErrorLogAll)
	}

	if err := s.updateTreeClusterPosition(ctx, id); err != nil {
		log.Error("error while update the cluster locations", "error", err, "cluster_id", id)
		return nil, service.MapError(ctx, err, service.ErrorLogAll)
	}

	log.Info("tree cluster updated successfully", "cluster_id", id)
	if err := s.publishUpdateEvent(ctx, prevTc); err != nil {
		return nil, service.MapError(ctx, err, service.ErrorLogAll)
	}

	return s.GetByID(ctx, id)
}

func (s *TreeClusterService) Delete(ctx context.Context, id int32) error {
	log := logger.GetLogger(ctx)
	_, err := s.treeClusterRepo.GetByID(ctx, id)
	if err != nil {
		return service.MapError(ctx, err, service.ErrorLogEntityNotFound)
	}

	if err := s.treeRepo.UnlinkTreeClusterID(ctx, id); err != nil {
		log.Debug("failed to unlink tree from tree cluster", "cluster_id", id, "error", err)
		return service.MapError(ctx, err, service.ErrorLogAll)
	}

	if err := s.treeClusterRepo.Delete(ctx, id); err != nil {
		log.Debug("failed to delete tree cluster", "error", err, "cluster_id", id)
		return service.MapError(ctx, err, service.ErrorLogAll)
	}

	log.Info("tree cluster deleted successfully", "cluster_id", id)
	return nil
}

func (s *TreeClusterService) UpdateWateringStatuses(ctx context.Context) error {
	log := logger.GetLogger(ctx)
	treeClusters, _, err := s.treeClusterRepo.GetAll(ctx, "")
	if err != nil {
		log.Error("failed to fetch tree cluster", "error", err)
		return err
	}

	cutoffTime := time.Now().Add(-24 * time.Hour) // 1 day ago
	for _, cluster := range treeClusters {
		// Do nothing if watering status is not »just watered«
		if cluster.WateringStatus != domain.WateringStatusJustWatered {
			continue
		}

		if cluster.LastWatered.Before(cutoffTime) {
			wateringStatus, err := s.getWateringStatusOfTreeCluster(ctx, cluster.ID)
			if err != nil {
				log.Error("failed to get watering status of cluster", "cluster_id", cluster.ID, "error", err)
				return err
			}

			err = s.treeClusterRepo.Update(ctx, cluster.ID, func(tc *domain.TreeCluster) (bool, error) {
				tc.WateringStatus = wateringStatus
				return true, nil
			})

			if err != nil {
				log.Error("failed to update watering status of tree cluster", "cluster_id", cluster.ID, "error", err)
			} else {
				log.Debug("watering status of tree cluster is updated by scheduler", "cluster_id", cluster.ID)
			}
		}
	}

	log.Info("watering status update for tree clusters completed successfully")
	return nil
}

func (s *TreeClusterService) Ready() bool {
	return s.treeClusterRepo != nil
}

// Update the tree cluster only after the trees have been updated to the database,
// otherwise the center point of the tree cluster cannot be set
func (s *TreeClusterService) updateTreeClusterPosition(ctx context.Context, id int32) error {
	log := logger.GetLogger(ctx)
	wateringStatus, err := s.getWateringStatusOfTreeCluster(ctx, id)
	if err != nil {
		log.Error("could not update watering status", "error", err)
	}

	err = s.treeClusterRepo.Update(ctx, id, func(tc *domain.TreeCluster) (bool, error) {
		lat, long, region, err := s.getUpdatedLatLong(ctx, tc)
		if err != nil {
			log.Debug("cancel transaction on updateting tree cluster position due to error", "error", err, "cluster_id", id)
			return false, err
		}

		if tc.Region != nil && tc.Region.ID != region.ID {
			tc.Region = region
			log.Debug("updating region in tree cluster position", "id", region.ID, "name", region.Name)
		}

		if tc.Latitude != lat || tc.Longitude != long {
			tc.Latitude = lat
			tc.Longitude = long
			tc.WateringStatus = wateringStatus

			log.Info("update tree cluster position due to changed trees inside the tree cluster", "cluster_id", id)
			log.Debug("detailed updated tree cluster position informations", "cluster_id", id,
				slog.Group("new_position", "latitude", *lat, "longitude", *long),
			)
		}

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
func (s *TreeClusterService) handlePrevTreeLocation(ctx context.Context, trees []*domain.Tree, updateFn func(context.Context, int32, func(tc *domain.TreeCluster) (bool, error)) error) error {
	log := logger.GetLogger(ctx)
	visitedClusters := make(map[int32]bool)
	for _, tree := range trees {
		if tree.TreeCluster == nil || tree.TreeCluster.ID == 0 {
			continue
		}

		if _, ok := visitedClusters[tree.TreeCluster.ID]; ok {
			continue
		}

		updateFunc := func(tc *domain.TreeCluster) (bool, error) {
			lat, long, region, err := s.getUpdatedLatLong(ctx, tc)
			if err != nil {
				return false, err
			}

			tc.Latitude = lat
			tc.Longitude = long
			tc.Region = region

			return true, nil
		}

		if err := updateFn(ctx, tree.TreeCluster.ID, updateFunc); err != nil {
			log.Error("failed to update tree cluster after handling prev tree locations", "error", err, "cluster_id", tree.TreeCluster.ID, "tree_id", tree.ID)
			return err
		}

		if err := s.publishUpdateEvent(ctx, tree.TreeCluster); err != nil {
			return err
		}

		visitedClusters[tree.TreeCluster.ID] = true
	}

	log.Info("successfully updated tree cluster locations from prev trees",
		"tree_ids", utils.Map(trees, func(t *domain.Tree) int32 { return t.ID }),
		"updated_clusters", utils.MapKeysSlice(visitedClusters, func(k int32, _ bool) int32 { return k }),
	)
	return nil
}

func (s *TreeClusterService) getTrees(ctx context.Context, ids []*int32) ([]*domain.Tree, error) {
	treeIDs := make([]int32, len(ids))
	for i, id := range ids {
		treeIDs[i] = *id
	}

	return s.treeRepo.GetTreesByIDs(ctx, treeIDs)
}

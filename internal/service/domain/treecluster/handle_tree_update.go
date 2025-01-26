package treecluster

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
)

func (s *TreeClusterService) HandleCreateTree(ctx context.Context, event *entities.EventCreateTree) error {
	log := logger.GetLogger(ctx)
	log.Debug("handle event", "event", event.Type(), "service", "TreeClusterService")

	if event.New.TreeCluster == nil {
		return nil
	}

	return s.handleTreeClusterUpdate(ctx, event.New.TreeCluster, event.New)
}

func (s *TreeClusterService) HandleDeleteTree(ctx context.Context, event *entities.EventDeleteTree) error {
	log := logger.GetLogger(ctx)
	log.Debug("handle event", "event", event.Type(), "service", "TreeClusterService")

	if event.Prev.TreeCluster == nil {
		return nil
	}

	return s.handleTreeClusterUpdate(ctx, event.Prev.TreeCluster, event.Prev)
}

func (s *TreeClusterService) HandleUpdateTree(ctx context.Context, event *entities.EventUpdateTree) error {
	log := logger.GetLogger(ctx)
	log.Debug("handle event", "event", event.Type(), "service", "TreeClusterService")

	if event.Prev.TreeCluster == nil && event.New.TreeCluster == nil {
		return nil
	}

	if s.isNoUpdateNeeded(event) {
		return nil
	}

	if err := s.handleTreeClusterUpdate(ctx, event.Prev.TreeCluster, event.New); err != nil {
		return err
	}

	if event.Prev.TreeCluster != nil && event.New.TreeCluster != nil && event.Prev.TreeCluster.ID != event.New.TreeCluster.ID {
		if err := s.handleTreeClusterUpdate(ctx, event.New.TreeCluster, event.New); err != nil {
			return err
		}
	}

	return nil
}

func (s *TreeClusterService) isNoUpdateNeeded(event *entities.EventUpdateTree) bool {
	treePosSame := event.Prev.Latitude == event.New.Latitude && event.Prev.Longitude == event.New.Longitude
	tcSame := event.Prev.TreeCluster != nil && event.New.TreeCluster != nil && event.Prev.TreeCluster.ID == event.New.TreeCluster.ID
	return treePosSame && tcSame
}

func (s *TreeClusterService) handleTreeClusterUpdate(ctx context.Context, tc *entities.TreeCluster, tree *entities.Tree) error {
	log := logger.GetLogger(ctx)
	if tc == nil {
		return nil
	}

	wateringStatus, err := s.getWateringStatusOfTreeCluster(ctx, tree)
	if err != nil {
		return nil
	}

	updateFn := func(tc *entities.TreeCluster) (bool, error) {
		lat, long, region, err := s.getUpdatedLatLong(ctx, tc)
		if err != nil {
			log.Error("failed to calculate latitude, longitude and region based on tree cluster", "error", err, "cluster_id", tc.ID)
			return false, err
		}
		tc.Latitude = lat
		tc.Longitude = long
		tc.Region = region
		tc.WateringStatus = wateringStatus
		return true, nil
	}

	if err := s.treeClusterRepo.Update(ctx, tc.ID, updateFn); err == nil {
		log.Info("successfully updated new tree cluster position", "cluster_id", tc.ID)
		return s.publishUpdateEvent(ctx, tc)
	}

	return nil
}

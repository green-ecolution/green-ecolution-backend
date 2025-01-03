package treecluster

import (
	"context"
	"log/slog"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
)

func (s *TreeClusterService) HandleUpdateTree(ctx context.Context, event *entities.EventUpdateTree) error {
	slog.Debug("handle event", "event", event.Type(), "service", "TreeClusterService")

	if event.Prev.TreeCluster == nil && event.New.TreeCluster == nil {
		return nil
	}

	if s.isNoUpdateNeeded(event) {
		return nil
	}

	if err := s.handleTreeClusterUpdate(ctx, event.Prev.TreeCluster); err != nil {
		return err
	}

	if event.Prev.TreeCluster.ID != event.New.TreeCluster.ID {
		if err := s.handleTreeClusterUpdate(ctx, event.New.TreeCluster); err != nil {
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

func (s *TreeClusterService) handleTreeClusterUpdate(ctx context.Context, tc *entities.TreeCluster) error {
	if tc == nil {
		return nil
	}

	updateFn := func(tc *entities.TreeCluster) (bool, error) {
		lat, long, region, err := s.getUpdatedLatLong(ctx, tc)
		if err != nil {
			return false, err
		}
		tc.Latitude = lat
		tc.Longitude = long
		tc.Region = region
		return true, nil
	}

	if err := s.treeClusterRepo.Update(ctx, tc.ID, updateFn); err == nil {
		return s.publishUpdateEvent(ctx, tc)
	}
	return nil
}

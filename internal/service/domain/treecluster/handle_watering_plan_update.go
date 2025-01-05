package treecluster

import (
	"context"
	"log/slog"
	"reflect"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
)

func (s *TreeClusterService) HandleUpdateWateringPlan(ctx context.Context, event *entities.EventUpdateWateringPlan) error {
	slog.Debug("handle event", "event", event.Type(), "service", "TreeClusterService")

	// Tree clusters should only be updated if the status has been changed to ‘finished’
	// and the linked tree clusters and the date have not changed
	if event.Prev.Status == event.New.Status ||
		event.Prev.Date != event.New.Date ||
		event.New.Status != entities.WateringPlanStatusFinished ||
		!reflect.DeepEqual(event.Prev.TreeClusters, event.New.TreeClusters) {
		return nil
	}

	if err := s.handleTreeClustersUpdate(ctx, event.New.TreeClusters, event.New.Date); err != nil {
		return err
	}

	return nil
}

func (s *TreeClusterService) handleTreeClustersUpdate(ctx context.Context, tcs []*entities.TreeCluster, date time.Time) error {
	if len(tcs) == 0 || tcs == nil {
		return nil
	}

	for _, tc := range tcs {
		updateFn := func(tc *entities.TreeCluster) (bool, error) {
			tc.LastWatered = &date
			return true, nil
		}

		if err := s.treeClusterRepo.Update(ctx, tc.ID, updateFn); err != nil {
			slog.Error("Failed to update tree cluster in response to event", "tree_cluster_id", tc.ID, "service", "TreeClusterService", "error", err)
		}
	}

	return nil
}

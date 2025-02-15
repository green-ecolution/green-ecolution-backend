package treecluster

import (
	"context"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
)

func (s *TreeClusterService) HandleUpdateWateringPlan(ctx context.Context, event *entities.EventUpdateWateringPlan) error {
	log := logger.GetLogger(ctx)
	log.Debug("handle event", "event", event.Type(), "service", "TreeClusterService")

	// Tree clusters should only be updated if the status has been changed to ‘finished’
	// and the linked tree clusters and the date have not changed
	if event.Prev.Status == event.New.Status ||
		event.Prev.Date != event.New.Date ||
		event.New.Status != entities.WateringPlanStatusFinished ||
		len(event.Prev.TreeClusters) != len(event.New.TreeClusters) {
		return nil
	}

	if err := s.handleTreeClustersUpdate(ctx, event.New.TreeClusters, event.New.Date); err != nil {
		return err
	}

	return nil
}

func (s *TreeClusterService) handleTreeClustersUpdate(ctx context.Context, tcs []*entities.TreeCluster, date time.Time) error {
	log := logger.GetLogger(ctx)
	if len(tcs) == 0 || tcs == nil {
		return nil
	}

	for _, tc := range tcs {
		updateFn := func(tc *entities.TreeCluster) (bool, error) {
			tc.LastWatered = &date
			return true, nil
		}

		if err := s.treeClusterRepo.Update(ctx, tc.ID, updateFn); err == nil {
			log.Info("successfully updated last watered date in tree cluster", "cluster_id", tc.ID, "last_watered", date)
			err := s.publishUpdateEvent(ctx, tc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

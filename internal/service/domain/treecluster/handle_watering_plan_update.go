package treecluster

import (
	"context"
	"log/slog"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
)

func (s *TreeClusterService) HandleUpdateWateringPlan(ctx context.Context, event *entities.EventUpdateWateringPlan) error {
	slog.Debug("handle event", "event", event.Type(), "service", "TreeClusterService")

	return nil
}

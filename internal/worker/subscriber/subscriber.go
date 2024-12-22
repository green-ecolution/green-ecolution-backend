package subscriber

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

type UpdateTreeSubscriber struct {
	tcs service.TreeClusterService
}

func NewUpdateTreeSubscriber(tcs service.TreeClusterService) *UpdateTreeSubscriber {
	return &UpdateTreeSubscriber{
		tcs: tcs,
	}
}

func (s *UpdateTreeSubscriber) EventType() entities.EventType {
	return entities.EventTypeUpdateTree
}

func (s *UpdateTreeSubscriber) HandleEvent(ctx context.Context, e entities.Event) error {
	event := e.(entities.EventUpdateTree)
	return s.tcs.HandleUpdateTree(ctx, event)
}

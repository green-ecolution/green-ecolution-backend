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
	return s.tcs.HandleUpdateTree(ctx, &event)
}

type CreateTreeSubscriber struct {
	tcs service.TreeClusterService
}

func NewCreateTreeSubscriber(tcs service.TreeClusterService) *CreateTreeSubscriber {
	return &CreateTreeSubscriber{
		tcs: tcs,
	}
}

func (s *CreateTreeSubscriber) EventType() entities.EventType {
	return entities.EventTypeUpdateTree
}

func (s *CreateTreeSubscriber) HandleEvent(ctx context.Context, e entities.Event) error {
	event := e.(entities.EventUpdateTree)
	return s.tcs.HandleUpdateTree(ctx, &event)
}

type DeleteTreeSubscriber struct {
	tcs service.TreeClusterService
}

func NewDeleteTreeSubscriber(tcs service.TreeClusterService) *DeleteTreeSubscriber {
	return &DeleteTreeSubscriber{
		tcs: tcs,
	}
}

func (s *DeleteTreeSubscriber) EventType() entities.EventType {
	return entities.EventTypeUpdateTree
}

func (s *DeleteTreeSubscriber) HandleEvent(ctx context.Context, e entities.Event) error {
	event := e.(entities.EventUpdateTree)
	return s.tcs.HandleUpdateTree(ctx, &event)
}

type CreateSensorDataSubscriber struct {
	treeSvc service.TreeService
}

func NewSensorDataSubscriber(tcSvc service.TreeClusterService, treeSvc service.TreeService) *CreateSensorDataSubscriber {
	return &CreateSensorDataSubscriber{
		treeSvc: treeSvc,
	}
}

func (s *CreateSensorDataSubscriber) EventType() entities.EventType {
	return entities.EventTypeNewSensorData
}

func (s *CreateSensorDataSubscriber) HandleEvent(ctx context.Context, e entities.Event) error {
	event := e.(entities.EventNewSensorData)
	return s.treeSvc.HandleNewSensorData(ctx, &event)
}

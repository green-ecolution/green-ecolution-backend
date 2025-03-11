package plugin

import (
	"context"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

var _ service.PluginService = (*DummyPluginManager)(nil)

type DummyPluginManager struct{}

func NewDummyPluginManager() *DummyPluginManager {
	return &DummyPluginManager{}
}

func (s *DummyPluginManager) Register(ctx context.Context, plugin *entities.Plugin) (*entities.ClientToken, error) {
	return nil, service.NewError(service.Gone, "plugin support is disabled")
}

func (s *DummyPluginManager) RefreshToken(ctx context.Context, auth *entities.AuthPlugin, slug string) (*entities.ClientToken, error) {
	return nil, service.NewError(service.Gone, "plugin support is disabled")
}

func (s *DummyPluginManager) Get(ctx context.Context, slug string) (entities.Plugin, error) {
	return entities.Plugin{}, service.NewError(service.NotFound, "plugin support is disabled")
}

func (s *DummyPluginManager) GetAll(ctx context.Context) ([]entities.Plugin, []time.Time) {
	return []entities.Plugin{}, []time.Time{}
}

func (s *DummyPluginManager) HeartBeat(ctx context.Context, slug string) error {
	return service.NewError(service.Gone, "plugin support is disabled")
}

func (s *DummyPluginManager) Unregister(ctx context.Context, slug string) {}

func (s *DummyPluginManager) StartCleanup(ctx context.Context) {}

func (s *DummyPluginManager) Ready() bool {
	return true
}

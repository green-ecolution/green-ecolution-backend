package plugin

import (
	"context"
	"log/slog"
	"time"
)

type PluginWorkerConfig struct {
	interval time.Duration
	timeout  time.Duration
}

type PluginWorker struct {
	cfg PluginWorkerConfig
}

type PluginWorkerOption func(*PluginWorkerConfig)

var defaultPluginWorkerConfig = PluginWorkerConfig{
	timeout:  5 * time.Minute,
	interval: 1 * time.Minute,
}

func WithTimeout(timeout time.Duration) PluginWorkerOption {
	return func(cfg *PluginWorkerConfig) {
		cfg.timeout = timeout
	}
}

func WithInterval(interval time.Duration) PluginWorkerOption {
	return func(cfg *PluginWorkerConfig) {
		cfg.interval = interval
	}
}

func NewPluginWorker(opts ...PluginWorkerOption) *PluginWorker {
	cfg := defaultPluginWorkerConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	return &PluginWorker{
		cfg: cfg,
	}
}

func (w *PluginWorker) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.cfg.interval)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			for name, plugin := range registeredPlugins {
				if time.Since(plugin.LastHeartbeat) > w.cfg.timeout {
					slog.Info("Removing plugin due to timeout", "plugin", name)
					delete(registeredPlugins, name)
				}
			}
		}
	}
}

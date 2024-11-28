package plugin

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"
)

// PluginWorkerConfig is the configuration for a PluginWorker.
type PluginWorkerConfig struct {
	plugin         Plugin
	host           *url.URL
	hostAPIVersion string
	interval       time.Duration
	client         *http.Client
	token          *Token
}

// PluginWorker is a worker that registers a plugin with the plugin host and sends heartbeats.
type PluginWorker struct {
	cfg PluginWorkerConfig
}

// PluginWorkerOption is a functional option for configuring a PluginWorker.
type PluginWorkerOption func(*PluginWorkerConfig)

// WithClient sets the client for the PluginWorker.
func WithClient(client *http.Client) PluginWorkerOption {
	return func(cfg *PluginWorkerConfig) {
		cfg.client = client
	}
}

// WithHost sets the host for the PluginWorker.
func WithHost(host *url.URL) PluginWorkerOption {
	return func(cfg *PluginWorkerConfig) {
		cfg.host = host
	}
}

// WithHostAPIVersion sets the host API version for the PluginWorker.
func WithHostAPIVersion(version string) PluginWorkerOption {
	return func(cfg *PluginWorkerConfig) {
		cfg.hostAPIVersion = version
	}
}

// WithPlugin sets the plugin for the PluginWorker.
func WithPlugin(plugin Plugin) PluginWorkerOption {
	return func(cfg *PluginWorkerConfig) {
		cfg.plugin = plugin
	}
}

// WithInterval sets the interval for the PluginWorker.
func WithInterval(interval time.Duration) PluginWorkerOption {
	return func(cfg *PluginWorkerConfig) {
		cfg.interval = interval
	}
}

// WithToken sets the token for the PluginWorker.
func WithToken(token *Token) PluginWorkerOption {
	return func(cfg *PluginWorkerConfig) {
		cfg.token = token
	}
}

// IsValid returns true if the PluginWorkerConfig is valid.
func (c *PluginWorkerConfig) IsValid() bool {
	return c.host != nil && c.plugin.PluginHostPath != nil && c.interval > 0 && c.client != nil && c.plugin.Name != ""
}

var defaultCfg = PluginWorkerConfig{
	client:         http.DefaultClient,
	interval:       2 * time.Minute,
	hostAPIVersion: "v1",
}

// NewPluginWorker creates a new PluginWorker with the provided options. If no options are provided, the default values are used.
func NewPluginWorker(opts ...PluginWorkerOption) (*PluginWorker, error) {
	cfg := defaultCfg
	for _, opt := range opts {
		opt(&cfg)
	}
	if !cfg.IsValid() {
		return nil, errors.New("invalid config")
	}

	return &PluginWorker{cfg: cfg}, nil
}

// RunHeartbeat runs the heartbeat runner for the PluginWorker.
func (w *PluginWorker) RunHeartbeat(ctx context.Context) error {
	ticker := time.NewTicker(w.cfg.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := w.Heartbeat(ctx); err != nil {
				return err
			}
		}
	}
}

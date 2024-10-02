package demoplugin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type PluginWorkerConfig struct {
	client     *http.Client
	host       *url.URL
	pluginName string
	pluginPath *url.URL
	interval   time.Duration
}

type PluginWorker struct {
	cfg PluginWorkerConfig
}

type PluginWorkerOption func(*PluginWorkerConfig)

func WithClient(client *http.Client) PluginWorkerOption {
	return func(cfg *PluginWorkerConfig) {
		cfg.client = client
	}
}

func WithHost(host *url.URL) PluginWorkerOption {
	return func(cfg *PluginWorkerConfig) {
		cfg.host = host
	}
}

func WithPluginName(pluginName string) PluginWorkerOption {
	return func(cfg *PluginWorkerConfig) {
		cfg.pluginName = pluginName
	}
}

func WithPluginPath(pluginPath *url.URL) PluginWorkerOption {
	return func(cfg *PluginWorkerConfig) {
		cfg.pluginPath = pluginPath
	}
}

func WithInterval(interval time.Duration) PluginWorkerOption {
	return func(cfg *PluginWorkerConfig) {
		cfg.interval = interval
	}
}

var defaultCfg = PluginWorkerConfig{
	client:   http.DefaultClient,
	interval: 500 * time.Millisecond,
}

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

func (c *PluginWorkerConfig) IsValid() bool {
	return c.client != nil && c.host != nil && c.pluginName != "" && c.pluginPath != nil
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name string      `json:"name"`
	Path string      `json:"path"`
	Auth AuthRequest `json:"auth"`
}

func (w *PluginWorker) Register(ctx context.Context, username, password string) error {
	reqBody := RegisterRequest{
		Name: w.cfg.pluginName,
		Path: w.cfg.pluginPath.String(),
    Auth: AuthRequest{
      Username: username,
      Password: password,
    },
	}

	buf, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	registerPath := fmt.Sprintf("%s://%s/api/v1/plugin/register", w.cfg.host.Scheme, w.cfg.host.Host)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, registerPath, bytes.NewReader(buf))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := w.cfg.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to register plugin")
	}

	return nil
}

func (w *PluginWorker) Heartbeat(ctx context.Context) error {
	registerPath := fmt.Sprintf("%s://%s/api/v1/plugin/%s/heartbeat", w.cfg.host.Scheme, w.cfg.host.Host, w.cfg.pluginName)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, registerPath, nil)
	if err != nil {
		return err
	}

	resp, err := w.cfg.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to send heartbeat")
	}

	return nil
}

func (w *PluginWorker) Run(ctx context.Context) error {
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

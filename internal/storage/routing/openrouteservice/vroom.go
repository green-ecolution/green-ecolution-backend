package openrouteservice

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
)

type VroomClientConfig struct {
	url    *url.URL
	client *http.Client
}

type VroomClientOption func(*VroomClientConfig)

type VroomClient struct {
	cfg VroomClientConfig
}

func WithClient(client *http.Client) VroomClientOption {
	return func(cfg *VroomClientConfig) {
		cfg.client = client
	}
}

func WithHostURL(url *url.URL) VroomClientOption {
	return func(cfg *VroomClientConfig) {
		cfg.url = url
	}
}

var defaultCfg = VroomClientConfig{
	client: http.DefaultClient,
}

func NewVroomClient(opts ...VroomClientOption) VroomClient {
	cfg := defaultCfg
	for _, opt := range opts {
		opt(&cfg)
	}
	return VroomClient{
		cfg: cfg,
	}
}

func (v *VroomClient) Send(ctx context.Context, reqBody *VroomReq) (*VroomResponse, error) {
	buf, err := json.Marshal(reqBody)
	if err != nil {
		slog.Error("failed to marshal vroom req body", "error", err, "req_body", reqBody)
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, v.cfg.url.String(), bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := v.cfg.client.Do(req)
	if err != nil {
		slog.Error("failed to send request to vroom service", "error", err)
		return nil, err
	}

	var vroomResp VroomResponse
	if err := json.NewDecoder(resp.Body).Decode(&vroomResp); err != nil {
		slog.Error("failed to decode vroom response")
	}

	return &vroomResp, nil
}

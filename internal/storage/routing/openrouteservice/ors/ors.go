package ors

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

type OrsClientConfig struct {
	url    *url.URL
	client *http.Client
}

type OrsClientOption func(*OrsClientConfig)

type OrsClient struct {
	cfg OrsClientConfig
}

func WithClient(client *http.Client) OrsClientOption {
	return func(cfg *OrsClientConfig) {
		cfg.client = client
	}
}

func WithHostURL(url *url.URL) OrsClientOption {
	return func(cfg *OrsClientConfig) {
		cfg.url = url
	}
}

var defaultCfg = OrsClientConfig{
	client: http.DefaultClient,
}

func NewOrsClient(opts ...OrsClientOption) OrsClient {
	cfg := defaultCfg
	for _, opt := range opts {
		opt(&cfg)
	}
	return OrsClient{
		cfg: cfg,
	}
}

func (o *OrsClient) DirectionsGeoJson(ctx context.Context, profile string, reqBody OrsDirectionRequest) (*OrsGeoJsonResponse, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(reqBody); err != nil {
		slog.Error("failed to marshal ors req body", "error", err, "req_body", reqBody)
		return nil, err
	}

	path := fmt.Sprintf("%s/v2/directions/%s/geojson", o.cfg.url.String(), profile)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := o.cfg.client.Do(req)
	if err != nil {
		slog.Error("failed to send request to ors service", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err == nil {
			slog.Error("response from the ORS service with a not successful code", "status_code", resp.StatusCode, "body", body)
		} else {
			slog.Error("response from the ORS service with a not successful code", "status_code", resp.StatusCode)
		}
		return nil, errors.New("response not successful")
	}

	var orsGeoJson OrsGeoJsonResponse
	if err := json.NewDecoder(resp.Body).Decode(&orsGeoJson); err != nil {
		slog.Error("failed to decode ors response")
	}

	return &orsGeoJson, nil
}

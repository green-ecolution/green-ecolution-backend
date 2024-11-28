package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
)

// Register registers the plugin with the plugin host.
func (w *PluginWorker) Register(ctx context.Context, clientID, clientSecret string) (*Token, error) {
	reqBody := entities.PluginRegisterRequest{
		Slug:        w.cfg.plugin.Slug,
		Name:        w.cfg.plugin.Name,
		Path:        w.cfg.plugin.PluginHostPath.String(),
		Version:     w.cfg.plugin.Version,
		Description: w.cfg.plugin.Description,
		Auth: entities.PluginAuth{
			ClientID:     clientID,
			ClientSecret: clientSecret,
		},
	}

	buf, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	registerPath := fmt.Sprintf("%s://%s/api/%s/plugin/register", w.cfg.host.Scheme, w.cfg.hostAPIVersion, w.cfg.host.Host)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, registerPath, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := w.cfg.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to register plugin")
	}

	var tokenResp entities.ClientTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	token := &Token{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
	}

	w.cfg.token = token
	return token, nil
}

// Heartbeat sends a heartbeat to the plugin host.
func (w *PluginWorker) Heartbeat(ctx context.Context) error {
	heartbeatPath := fmt.Sprintf("%s://%s/api/%s/plugin/%s/heartbeat", w.cfg.host.Scheme, w.cfg.host.Host, w.cfg.hostAPIVersion, w.cfg.plugin.Slug)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, heartbeatPath, http.NoBody)
	if err != nil {
		return err
	}

	resp, err := w.cfg.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to send heartbeat")
	}

	return nil
}

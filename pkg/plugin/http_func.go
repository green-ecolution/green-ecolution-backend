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

// Register registers the plugin with the plugin host and returns an authentication token.
//
// This function performs the following steps:
//  1. Constructs a registration request payload containing plugin metadata (such as slug, name, version, and description)
//     and authentication credentials (client ID and client secret).
//  2. Sends the registration request as an HTTP POST request to the plugin host's registration API endpoint.
//  3. Parses the response from the plugin host to extract and return the authentication tokens.
//
// Parameters:
// - ctx: The context for managing request deadlines and cancellation.
// - clientID: The client identifier used for authenticating the plugin.
// - clientSecret: The secret key corresponding to the client identifier.
//
// Returns:
//   - A pointer to a Token struct containing the access and refresh tokens if registration is successful.
//   - An error if the registration fails, due to either an HTTP request issue, an unexpected status code,
//     or failure in parsing the response body.
//
// Possible errors:
// - If the request payload cannot be marshaled into JSON, an error is returned.
// - If creating the HTTP request fails, an error is returned.
// - If the HTTP request fails (e.g., network error), an error is returned.
// - If the plugin host returns a non-200 status code, an error is returned with a generic failure message.
// - If the response body cannot be decoded into the expected format, an error is returned.
//
// Example usage:
//
//	token, err := pluginWorker.Register(ctx, "client-id", "client-secret")
//	if err != nil {
//		log.Fatalf("Failed to register plugin: %v", err)
//	}
//	log.Printf("Successfully registered. Access token: %s", token.AccessToken)
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

	registerPath := fmt.Sprintf("%s://%s/api/%s/plugin/register", w.cfg.host.Scheme, w.cfg.host.Host, w.cfg.hostAPIVersion)
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

// Heartbeat sends a periodic heartbeat signal to the plugin host to indicate that the plugin is active and operational.
//
// This function performs the following steps:
// 1. Constructs the heartbeat API endpoint URL based on the plugin's slug and configuration settings.
// 2. Sends an HTTP POST request to the constructed endpoint without a request body.
// 3. Checks the response status code to confirm the success of the heartbeat operation.
//
// Parameters:
// - ctx: The context for managing request deadlines and cancellation.
//
// Returns:
// - nil if the heartbeat is successfully sent and acknowledged by the plugin host.
// - An error if any of the following issues occur:
//   - Constructing the HTTP request fails.
//   - Sending the HTTP request fails (e.g., due to network issues).
//   - The plugin host responds with a non-200 status code, indicating the heartbeat was not successfully received.
//
// Example usage:
//
//	err := pluginWorker.Heartbeat(ctx)
//	if err != nil {
//		log.Printf("Failed to send heartbeat: %v", err)
//	} else {
//		log.Println("Heartbeat sent successfully.")
//	}
//
// Notes:
// - The plugin's slug, host, and API version must be correctly configured in the `PluginWorker` instance.
// - This function assumes that the plugin host's heartbeat API endpoint requires an HTTP POST request.
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

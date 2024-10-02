package plugin

import (
	"fmt"
	"log/slog"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

func getPluginFiles(c *fiber.Ctx) error {
	pluginMutex.RLock()
	plugin := registeredPlugins[c.Params("plugin")]
	pluginMutex.RUnlock()

	reverseProxy := httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(&plugin.Path)
			r.Out.Host = r.In.Host
			r.Out.URL.Path = strings.Replace(r.In.URL.Path, "/api/v1/plugin/"+plugin.Name, plugin.Path.String(), 1)
			r.SetXForwarded()
		},
	}

	return adaptor.HTTPHandler(&reverseProxy)(c)
}

type PluginRegisterRequest struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Auth struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"auth"`
}

type PluginRegisterResponse struct {
	Success bool   `json:"success"`
	Token   entities.ClientTokenResponse `json:"token"`
}

func registerPlugin(svc service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req PluginRegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

    if req.Auth.Username == "" || req.Auth.Password == "" {
      return c.Status(fiber.StatusBadRequest).SendString("Username or password is empty")
    } 

    // Authenticate the plugin
    token, err := svc.AuthPlugin(c.Context(), &domain.AuthPlugin{
      Username: req.Auth.Username,
      Password: req.Auth.Password,
    })
    if err != nil {
      return err
    }

		path, err := url.Parse(req.Path)
		if err != nil {
			slog.Error("Failed to parse plugin path", "error", err)
			return c.Status(fiber.StatusBadRequest).SendString("Failed to parse plugin path")
		}

		pluginMutex.Lock()
		defer pluginMutex.Unlock()

		if _, ok := registeredPlugins[req.Name]; ok {
			return c.Status(fiber.StatusForbidden).SendString("plugin already registered")
		}

		registeredPlugins[req.Name] = Plugin{
			Name:          req.Name,
			Path:          *path,
			LastHeartbeat: time.Now(),
		}

		slog.Info("Plugin registered", "plugin", req.Name)
		slog.Debug("Plugin registered", "plugin", fmt.Sprintf("%+v", registeredPlugins[req.Name]))

		response := entities.ClientTokenResponse{
			AccessToken:  token.AccessToken,
			ExpiresIn:    token.ExpiresIn,
			RefreshToken: token.RefreshToken,
			TokenType:    token.TokenType,
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}

func pluginHeartbeat() fiber.Handler {
	return func(c *fiber.Ctx) error {
		pluginMutex.Lock()
		defer pluginMutex.Unlock()

		plugin := registeredPlugins[c.Params("plugin")]
		plugin.LastHeartbeat = time.Now()
		registeredPlugins[c.Params("plugin")] = plugin

		return c.Status(fiber.StatusOK).SendString("Heartbeat received")
	}
}

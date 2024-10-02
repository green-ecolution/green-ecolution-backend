package plugin

import (
	"log/slog"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func getPluginFiles(c *fiber.Ctx) error {
	pluginMutex.RLock()
	plugin := registeredPlugins[c.Params("plugin")]
  pluginMutex.RUnlock()

	reverseProxy := httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(plugin.Path)
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
}

func registerPlugin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req PluginRegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
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
			Path:          path,
			LastHeartbeat: time.Now(),
		}

		slog.Info("Plugin registered", "plugin", req.Name)
		return c.Status(fiber.StatusOK).SendString("Plugin registered")
	}
}

func pluginHeartbeat() fiber.Handler {
	return func(c *fiber.Ctx) error {
		pluginMutex.Lock()
		defer pluginMutex.Unlock()

		registeredPlugins[c.Params("plugin")] = Plugin{
			LastHeartbeat: time.Now(),
		}

		return c.Status(fiber.StatusOK).SendString("Heartbeat received")
	}
}

package plugin

import (
	"log/slog"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/spf13/viper"
)

func getPluginFiles(c *fiber.Ctx) error {
	url, err := url.Parse(viper.GetString("server.app_url"))
	if err != nil {
		return err
	}

	plugin := c.Params("plugin")

	reverseProxy := httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(url)
			r.Out.Host = r.In.Host
			r.Out.URL.Path = strings.Replace(r.In.URL.Path, "/api/v1/plugin/"+plugin, "/api/v1/demo_plugin", 1)
			r.SetXForwarded()
		},
	}

	return adaptor.HTTPHandler(&reverseProxy)(c)
}

type PluginRegisterRequest struct {
	Name string `json:"name"`
	Host string `json:"host"`
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

		pluginMutex.Lock()
		defer pluginMutex.Unlock()

		if _, ok := registeredPlugins[req.Name]; ok {
			return c.Status(fiber.StatusForbidden).SendString("plugin already registered")
		}

		registeredPlugins[req.Name] = Plugin{
			Name:          req.Name,
			Path:          req.Path,
			Host:          req.Host,
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

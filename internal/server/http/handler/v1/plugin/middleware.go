package plugin

import (
	"net/url"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Plugin struct {
	Name          string
	Path          url.URL
	LastHeartbeat time.Time
}

var pluginMutex sync.RWMutex
var registeredPlugins = make(map[string]Plugin)

func newPluginMiddleware() fiber.Handler {

	return func(c *fiber.Ctx) error {
		pluginMutex.RLock()
		if _, ok := registeredPlugins[c.Params("plugin")]; !ok {
			pluginMutex.RUnlock()
			return c.Status(fiber.StatusNotFound).SendString("plugin not found")
		}

		pluginMutex.RUnlock()
		return c.Next()
	}
}

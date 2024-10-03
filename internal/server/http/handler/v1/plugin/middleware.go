package plugin

import (
	// "sync"

	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
)

// var pluginMutex sync.RWMutex
var registeredPlugins = make(map[string]entities.Plugin)

func newPluginMiddleware() fiber.Handler {

	return func(c *fiber.Ctx) error {
		// pluginMutex.RLock()
		if _, ok := registeredPlugins[c.Params("plugin")]; !ok {
			// pluginMutex.RUnlock()
			return c.Status(fiber.StatusNotFound).SendString("plugin not found")
		}

		// pluginMutex.RUnlock()
		return c.Next()
	}
}

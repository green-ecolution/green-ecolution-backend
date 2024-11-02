package plugin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
)

var registeredPlugins = make(map[string]*entities.Plugin)

func newPluginMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if _, ok := registeredPlugins[c.Params("plugin")]; !ok {
			return c.Status(fiber.StatusNotFound).SendString("plugin not found")
		}

		return c.Next()
	}
}

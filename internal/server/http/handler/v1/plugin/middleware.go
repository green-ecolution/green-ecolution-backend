package plugin

import (
	"slices"

	"github.com/gofiber/fiber/v2"
)

func newPluginMiddleware() fiber.Handler {
  pluginNames := []string{
		"demo_plugin",
	}

	return func(c *fiber.Ctx) error {
		if !slices.Contains(pluginNames, c.Params("plugin")) {
			return c.Status(fiber.StatusNotFound).SendString("plugin not found")
		}

		return c.Next()
	}
}


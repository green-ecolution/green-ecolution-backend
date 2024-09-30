package plugin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

func RegisterRoutes(svc service.RegionService) *fiber.App {
	app := fiber.New()

  grp := app.Group("/csv_plugin", newPluginMiddleware())
  grp.Mount("/", getPluginFiles())

	return app
}

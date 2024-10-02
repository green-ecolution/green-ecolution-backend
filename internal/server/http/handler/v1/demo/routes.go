package demo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

func RegisterRoutes(svc service.RegionService) *fiber.App {
	app := fiber.New()

  app.Mount("/", getPluginFiles())
	return app
}

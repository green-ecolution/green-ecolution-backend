package route

import (
	"github.com/SmartCityFlensburg/green-space-management/internal/service"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(svc service.RouteService) *fiber.App {
	app := fiber.New()

	app.Get("/", GetRouteAll(svc))
	return app
}

package region

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

func RegisterRoutes(svc service.RegionService) *fiber.App {
	app := fiber.New()

	app.Post("/", GetAllRegions(svc))
	app.Get("/:id", GetRegionByID(svc))

	return app
}

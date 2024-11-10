package vehicle

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

func RegisterRoutes(svc service.VehicleService) *fiber.App {
	app := fiber.New()

	app.Get("/", GetAllVehicles(svc))
	app.Get("/:id", GetVehicleByID(svc))
	app.Get("/plate/:plate", GetVehicleByPlate(svc))

	return app
}
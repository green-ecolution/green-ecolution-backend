package sensor

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

func RegisterRoutes(svc service.Service) *fiber.App {
	app := fiber.New()

	app.Get("/", GetAllSensor(svc))
	app.Get("/:id", GetSensorByID(svc))
	app.Get("/:id/data", GetSensorDataByID(svc))

	return app
}

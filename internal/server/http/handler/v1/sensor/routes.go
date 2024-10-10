package sensor

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

func RegisterRoutes(svc service.SensorService) *fiber.App {
	app := fiber.New()

	app.Get("/", GetAllSensors(svc))
	app.Get("/:id", GetSensorByID(svc))
	app.Get("/:id/data", GetSensorDataByID(svc))

	app.Post("/", CreateSensor(svc))
	app.Put("/:id", UpdateSensor(svc))
	app.Delete("/:id", DeleteSensor(svc))

	return app
}

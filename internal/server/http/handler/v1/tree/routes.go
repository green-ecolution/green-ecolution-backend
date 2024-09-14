package tree

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

func RegisterRoutes(svc service.TreeService) *fiber.App {
	app := fiber.New()

	app.Get("/", GetAllTrees(svc))
	app.Get("/:id", GetTreeByID(svc))
	app.Patch("/:id", UpdateTree(svc))
	app.Post("/", CreateTree(svc))
	app.Delete("/", DeleteTree(svc))

	app.Get("/:id/images", GetTreeImages(svc))
	app.Post("/:id/images", AddTreeImage(svc))
	app.Delete("/:id/images/:image_id", RemoveTreeImage(svc))

	app.Get("/:id/sensor", GetTreeSensor(svc))
	app.Post("/:id/sensor", AddTreeSensor(svc))
	app.Delete("/:id/sensor/:sensor_id", RemoveTreeSensor(svc))

	return app
}

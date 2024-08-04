package tree

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

func RegisterRoutes(svc service.TreeService) *fiber.App {
	app := fiber.New()

	app.Get("/", GetAllTree(svc))
	app.Get("/:id", GetTreeByID(svc))
	app.Get("/:id/prediction", GetTreePredictions(svc))

	return app
}

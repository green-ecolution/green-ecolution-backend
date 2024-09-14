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

	return app
}

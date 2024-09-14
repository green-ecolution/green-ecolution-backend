package role

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

func RegisterRoutes(svc service.AuthService) *fiber.App {
	app := fiber.New()

	app.Post("/", GetAllUserRoles(svc))
  app.Get("/:id", GetRoleByID(svc))

	return app
}

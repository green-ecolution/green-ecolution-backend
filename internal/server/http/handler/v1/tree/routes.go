package tree

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

func RegisterRoutes(svc service.InfoService) *fiber.App {
	app := fiber.New()

	app.Get("/", GetAppInfo(svc))

	return app
}

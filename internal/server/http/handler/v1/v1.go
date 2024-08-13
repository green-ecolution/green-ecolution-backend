package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/green-ecolution/green-ecolution-backend/docs"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/info"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/tree"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

func V1Handler(svc *service.Services) *fiber.App {
	app := fiber.New()

	app.Mount("/info", info.RegisterRoutes(svc.InfoService))
	app.Mount("/tree", tree.RegisterRoutes(svc.TreeService))
	app.Get("/swagger/*", swagger.HandlerDefault) // default

	return app
}

package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/info"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/tree"

	"github.com/gofiber/swagger"
	_ "github.com/green-ecolution/green-ecolution-backend/docs"
)

func (s *Server) router() *fiber.App {
	app := fiber.New()

	app.Mount("/info", info.RegisterRoutes(s.services.InfoService))
	app.Mount("/treeSQL", tree.RegisterRoutes(s.services.TreeService))

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	return app
}

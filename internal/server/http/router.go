package http

import (
	"github.com/SmartCityFlensburg/green-space-management/internal/server/http/handler/info"
	"github.com/SmartCityFlensburg/green-space-management/internal/server/http/handler/tree"
	"github.com/SmartCityFlensburg/green-space-management/internal/server/http/handler/route"
	"github.com/gofiber/fiber/v2"

	_ "github.com/SmartCityFlensburg/green-space-management/docs"
	"github.com/gofiber/swagger"
)

func (s *Server) router() *fiber.App {
	app := fiber.New()

	app.Mount("/info", info.RegisterRoutes(s.services.InfoService))
	app.Mount("/tree", tree.RegisterRoutes(s.services.TreeService))
  app.Mount("/route", route.RegisterRoutes(s.services.RouteService))

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	return app
}

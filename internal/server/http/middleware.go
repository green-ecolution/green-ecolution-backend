package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/middleware"
)

func (s *Server) middleware() *fiber.App {
	app := fiber.New()

	app.Use(middleware.HealthCheck(s.services))
  app.Use(middleware.HttpLogger())

	return app
}

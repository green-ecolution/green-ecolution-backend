package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	v1 "github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/user"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/middleware"
)

func (s *Server) router() *fiber.App {
	app := fiber.New()

	app.Mount("/v1", v1.V1Handler(s.services))

	return app
}

func (s *Server) privateRoutes(app *fiber.App) {
}

func (s *Server) publicRoutes(app *fiber.App) {
	app.Use("/", middleware.HealthCheck(s.services))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	grp := app.Group("/api/v1")
	grp.Get("/swagger/*", swagger.HandlerDefault)
	grp.Post("/user", user.Register(s.services.AuthService))
}

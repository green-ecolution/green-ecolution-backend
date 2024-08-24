package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/auth"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/info"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/tree"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/middleware"
)

func (s *Server) privateRoutes(app *fiber.App) {
	grp := app.Group("/api/v1")

	grp.Mount("/info", info.RegisterRoutes(s.services.InfoService))
	grp.Mount("/tree", tree.RegisterRoutes(s.services.TreeService))
}

func (s *Server) publicRoutes(app *fiber.App) {
	app.Use("/", middleware.HealthCheck(s.services))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	grp := app.Group("/api/v1")
	grp.Get("/swagger/*", swagger.HandlerDefault)
	grp.Post("/user", auth.Register(s.services.AuthService))
  grp.Get("/login", auth.Login(s.services.AuthService))
  grp.Post("/token", auth.RequestToken(s.services.AuthService))
}

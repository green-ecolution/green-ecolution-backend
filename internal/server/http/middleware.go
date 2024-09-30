package http

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/middleware"
)

func (s *Server) middleware(
	initPublicRoutes func(app *fiber.App),
	initPrivateRoutes func(app *fiber.App),
) *fiber.App {
	slog.Info("Setting up fiber middlewares")

	app := fiber.New()

	app.Use(middleware.HealthCheck(s.services))
	app.Use(middleware.HTTPLogger())
	app.Use(middleware.RequestID())

	initPublicRoutes(app)
	initPrivateRoutes(app)
	//app.Use(middleware.NewJWTMiddleware(&s.cfg.IdentityAuth, s.services.AuthService))

	slog.Info("Fiber middlewares setup complete")

	return app
}

package http

import (
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/middleware"
	"github.com/spf13/viper"
)

func (s *Server) middleware() *fiber.App {
	slog.Info("Setting up fiber middlewares")
	logFormat := viper.GetString("server.logs.format")
	logLevel := viper.GetString("server.logs.level")

	logFn := logger.CreateLogger(os.Stdout, logger.LogFormat(logFormat), logger.LogLevel(logLevel))

	app := fiber.New()

	app.Use(middleware.HealthCheck(s.services))
	app.Use(middleware.HTTPLogger())
	app.Use(middleware.RequestID())
	app.Use(middleware.AppLogger(logFn))

	authMiddlware := middleware.NewJWTMiddleware(&s.cfg.IdentityAuth, s.services.AuthService)
	s.root(app, authMiddlware)
	slog.Info("Fiber middlewares setup complete")

	return app
}

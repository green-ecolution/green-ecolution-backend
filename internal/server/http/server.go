package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/worker"

	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

type HTTPError struct {
	Error  string `json:"error"`
	Code   int    `json:"code"`
	Path   string `json:"path"`
	Method string `json:"method"`
} // @Name HTTPError

type Server struct {
	cfg      *config.Config
	services *service.Services
}

func NewServer(cfg *config.Config, services *service.Services) *Server {
	return &Server{
		cfg:      cfg,
		services: services,
	}
}

func (s *Server) Run(ctx context.Context) error {
	app := fiber.New(fiber.Config{
		AppName:      s.cfg.Dashboard.Title,
		ServerHeader: s.cfg.Dashboard.Title,
		ErrorHandler: errorHandler,
	})

	app.Mount("/", s.middleware())

	go func() {
		slog.Info("starting plugin cleanup service: cleaning up unhealthy plugins")
		err := s.services.PluginService.StartCleanup(ctx)
		if err != nil {
			slog.Error("error while running plugin cleanup", "error", err)
		}
	}()

	go func() {
		worker.Scheduler(ctx, 1*time.Hour, s.services.SensorService.UpdateStatuses)
	}()

	go func() {
		<-ctx.Done()
		slog.Info("shutting down http server")
		if err := app.Shutdown(); err != nil {
			slog.Error("error while shutting down http server", "error", err, "service", "fiber")
		}
	}()

	slog.Info("starting server", "url", s.cfg.Server.AppURL, "port", s.cfg.Server.Port, "service", "fiber")
	return app.Listen(fmt.Sprintf(":%d", s.cfg.Server.Port))
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	return c.Status(code).JSON(HTTPError{
		Error:  err.Error(),
		Code:   code,
		Path:   c.Path(),
		Method: c.Method(),
	})
}

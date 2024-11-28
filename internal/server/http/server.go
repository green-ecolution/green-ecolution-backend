package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

const (
	pluginCleanupTimeout = 5 * time.Second
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
		err := s.services.PluginService.StartCleanup(ctx)
		slog.Error("Error while running plugin cleanup", "error", err)
	}()

	go func() {
		<-ctx.Done()
		fmt.Println("Shutting down HTTP Server")
		if err := app.Shutdown(); err != nil {
			fmt.Println("Error while shutting down HTTP Server:", err)
		}
	}()

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

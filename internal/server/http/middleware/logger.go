package middleware

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	http_logger "github.com/gofiber/fiber/v2/middleware/logger"
)

func HTTPLogger() fiber.Handler {
	return http_logger.New()
}

func AppLogger(createLoggerFn func() *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := createLoggerFn()
		requestid := c.Locals("requestid").(string)

		logger = logger.With("request_id", requestid)
		c.Locals("logger", logger)

		return c.Next()
	}
}

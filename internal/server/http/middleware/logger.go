package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	http_logger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
)

func HTTPLogger() fiber.Handler {
	return http_logger.New()
}

func AppLogger(createLoggerFn func() *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log := createLoggerFn()
		requestid := c.Locals("requestid").(string)

		log = log.With("request_id", requestid, "request_duration", logger.NewTimeSince(), "request_start_time", time.Now())
		c.Locals("logger", log)

		return c.Next()
	}
}

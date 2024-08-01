package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func HTTPLogger() fiber.Handler {
	return logger.New()
}

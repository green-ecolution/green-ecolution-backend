package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
)

func (s *Server) healthCheck() func(c *fiber.Ctx) error {
	return healthcheck.New(healthcheck.Config{
		LivenessProbe: func(_ *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/health",
		ReadinessProbe: func(_ *fiber.Ctx) bool {
			return s.services.AllServicesReady()
		},
		ReadinessEndpoint: "/ready",
	})
}

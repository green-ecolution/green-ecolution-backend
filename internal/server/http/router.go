package http

import (
	"github.com/gofiber/fiber/v2"
	v1 "github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1"
)

func (s *Server) router() *fiber.App {
	app := fiber.New()

	app.Mount("/v1", v1.V1Handler(s.services))

	return app
}

package routing

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

func RegisterRoutes(r fiber.Router, svc service.RoutingService) {
	r.Post("/preview", CreatePreviewRoute(svc))
	r.Post("/", CreateRoute(svc))
}

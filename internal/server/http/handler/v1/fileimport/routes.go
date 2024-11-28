package fileimport

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

func RegisterRoutes(r fiber.Router, svc service.TreeService) {
	r.Post("/csv", ImportTreesFromCSV(svc))
}

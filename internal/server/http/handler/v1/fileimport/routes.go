package fileimport

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

func RegisterRoutes(svc service.TreeService) *fiber.App {
	router := fiber.New()
	router.Post("/csv", ImportTreesFromCSV(svc))
	return router
}

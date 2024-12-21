package vehicle

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

func RegisterRoutes(r fiber.Router, svc service.VehicleService) {
	r.Get("/", GetAllVehicles(svc))
	r.Get("/type/:type", GetAllVehiclesByType(svc))
	r.Get("/:id", GetVehicleByID(svc))
	r.Get("/plate/:plate", GetVehicleByPlate(svc))
	r.Post("/", CreateVehicle(svc))
	r.Put("/:id", UpdateVehicle(svc))
	r.Delete("/:id", DeleteVehicle(svc))
}

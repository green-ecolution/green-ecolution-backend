package tree

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

func RegisterRoutes(r fiber.Router, svc service.TreeService) {
	r.Get("/", GetAllTrees(svc))
	r.Get("/:id", GetTreeByID(svc))
	r.Put("/:id", UpdateTree(svc))
	r.Post("/", CreateTree(svc))
	r.Delete("/:id", DeleteTree(svc))

	r.Get("/:id/images", GetTreeImages(svc))
	r.Post("/:id/images", AddTreeImage(svc))
	r.Delete("/:id/images/:image_id", RemoveTreeImage(svc))

	r.Get("/:id/sensor", GetTreeSensor(svc))
	r.Post("/:id/sensor", AddTreeSensor(svc))
	r.Delete("/:id/sensor/:sensor_id", RemoveTreeSensor(svc))

	r.Get("/sensor/:sensor_id", GetTreeBySensorID(svc))
}

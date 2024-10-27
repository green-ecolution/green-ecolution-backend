package treecluster

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

func RegisterRoutes(svc service.TreeClusterService) *fiber.App {
	app := fiber.New()

	app.Get("/", GetAllTreeClusters(svc))
	app.Get("/:treecluster_id", GetTreeClusterByID(svc))
	app.Post("/", CreateTreeCluster(svc))
	app.Put("/:treecluster_id", UpdateTreeCluster(svc))
	app.Delete("/:treecluster_id", DeleteTreeCluster(svc))
	
	return app
}

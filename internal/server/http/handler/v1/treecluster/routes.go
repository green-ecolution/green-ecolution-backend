package treecluster

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

func RegisterRoutes(svc service.Service) *fiber.App {
	app := fiber.New()

	app.Get("/", GetAllTreeClusters(svc))
	app.Get("/:treecluster_id", GetTreeClusterByID(svc))
	app.Post("/", CreateTreeCluster(svc))
	app.Put("/:treecluster_id", UpdateTreeCluster(svc))
	app.Delete("/:treecluster_id", DeleteTreeCluster(svc))
	app.Get("/:treecluster_id/trees", GetTreesInTreeCluster(svc))
  app.Post("/:treecluster_id/trees", AddTreesToTreeCluster(svc))
  app.Delete("/:treecluster_id/trees/:tree_id", RemoveTreesFromTreeCluster(svc))

	return app
}

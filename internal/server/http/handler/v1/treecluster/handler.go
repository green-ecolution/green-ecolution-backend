package treecluster

import (
	"github.com/gofiber/fiber/v2"
	_ "github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

// @Summary		Get all tree clusters
// @Description	Get all tree clusters
// @Id				get-all-tree-clusters
// @Tags			Tree Cluster
// @Produce		json
// @Success		200	{object}	entities.TreeClusterListResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/cluster [get]
// @Param			page			query	string	false	"Page"
// @Param			limit			query	string	false	"Limit"
// @Param			status			query	string	false	"Status"
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func GetAllTreeClusters(_ service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

// @Summary		Get tree cluster by ID
// @Description	Get tree cluster by ID
// @Id				get-tree-cluster-by-id
// @Tags			Tree Cluster
// @Produce		json
// @Success		200	{object}	entities.TreeClusterResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/cluster/{cluster_id} [get]
// @Param			cluster_id		path	string	true	"Tree Cluster ID"
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func GetTreeClusterByID(_ service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

// @Summary		Create tree cluster
// @Description	Create tree cluster
// @Id				create-tree-cluster
// @Tags			Tree Cluster
// @Produce		json
// @Success		200	{object}	entities.TreeClusterResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/cluster [post]
// @Param			body			body	entities.TreeClusterCreateRequest	true	"Tree Cluster Create Request"
// @Param			Authorization	header	string								true	"Insert your access token"	default(Bearer <Add access token here>)
func CreateTreeCluster(_ service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

// @Summary		Update tree cluster
// @Description	Update tree cluster
// @Id				update-tree-cluster
// @Tags			Tree Cluster
// @Produce		json
// @Success		200	{object}	entities.TreeClusterResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/cluster/{cluster_id} [put]
// @Param			cluster_id		path	string								true	"Tree Cluster ID"
// @Param			body			body	entities.TreeClusterUpdateRequest	true	"Tree Cluster Update Request"
// @Param			Authorization	header	string								true	"Insert your access token"	default(Bearer <Add access token here>)
func UpdateTreeCluster(_ service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

// @Summary		Delete tree cluster
// @Description	Delete tree cluster
// @Id				delete-tree-cluster
// @Tags			Tree Cluster
// @Produce		json
// @Success		200	{object}	entities.TreeClusterResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/cluster/{cluster_id} [delete]
// @Param			cluster_id		path	string	true	"Tree Cluster ID"
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func DeleteTreeCluster(_ service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

// @Summary		Get all trees in tree cluster
// @Description	Get all trees in tree cluster
// @Id				get-all-trees-in-tree-cluster
// @Tags			Tree Cluster
// @Produce		json
// @Success		200	{object}	entities.TreeResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/cluster/{cluster_id}/trees [get]
// @Param			cluster_id		path	string	true	"Tree Cluster ID"
// @Param			page			query	string	false	"Page"
// @Param			limit			query	string	false	"Limit"
// @Param			age				query	string	false	"Age"
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func GetTreesInTreeCluster(_ service.Service) fiber.Handler {
	// TODO: Change response @Success to entities.TreeListResponse
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

// @Summary		Add trees to tree cluster
// @Description	Add trees to tree cluster
// @Id				add-trees-to-tree-cluster
// @Tags			Tree Cluster
// @Produce		json
// @Success		200	{object}	entities.TreeClusterResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/cluster/{cluster_id}/trees [post]
// @Param			cluster_id		path	string								true	"Tree Cluster ID"
// @Param			body			body	entities.TreeClusterAddTreesRequest	true	"Tree Cluster Add Trees Request"
// @Param			Authorization	header	string								true	"Insert your access token"	default(Bearer <Add access token here>)
func AddTreesToTreeCluster(_ service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

// @Summary		Remove trees from tree cluster
// @Description	Remove trees from tree cluster
// @Id				remove-trees-from-tree-cluster
// @Tags			Tree Cluster
// @Produce		json
// @Success		200	{object}	entities.TreeClusterResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/cluster/{cluster_id}/trees/{tree_id} [delete]
// @Param			cluster_id		path	string	true	"Tree Cluster ID"
// @Param			tree_id			path	string	true	"Tree ID"
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func RemoveTreesFromTreeCluster(_ service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

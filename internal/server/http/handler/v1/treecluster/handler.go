package treecluster

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/tree"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/treecluster"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

//	@Summary		Get all tree clusters
//	@Description	Get all tree clusters
//	@Id				get-all-tree-clusters
//	@Tags			Tree Cluster
//	@Produce		json
//	@Success		200	{object}	treecluster.TreeClusterListResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/cluster [get]
//	@Param			page			query	string	false	"Page"
//	@Param			limit			query	string	false	"Limit"
//	@Param			status			query	string	false	"Status"
//	@Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func GetAllTreeClusters(svc service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.JSON(treecluster.TreeClusterListResponse{})
	}
}

//	@Summary		Get tree cluster by ID
//	@Description	Get tree cluster by ID
//	@Id				get-tree-cluster-by-id
//	@Tags			Tree Cluster
//	@Produce		json
//	@Success		200	{object}	treecluster.TreeClusterResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/cluster/{treecluster_id} [get]
//	@Param			treecluster_id	path	string	true	"Tree Cluster ID"
//	@Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func GetTreeClusterByID(svc service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.JSON(treecluster.TreeClusterResponse{})
	}
}

//	@Summary		Create tree cluster
//	@Description	Create tree cluster
//	@Id				create-tree-cluster
//	@Tags			Tree Cluster
//	@Produce		json
//	@Success		200	{object}	treecluster.TreeClusterResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/cluster [post]
//	@Param			body			body	treecluster.TreeClusterCreateRequest	true	"Tree Cluster Create Request"
//	@Param			Authorization	header	string									true	"Insert your access token"	default(Bearer <Add access token here>)
func CreateTreeCluster(svc service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.JSON(treecluster.TreeClusterResponse{})
	}
}

//	@Summary		Update tree cluster
//	@Description	Update tree cluster
//	@Id				update-tree-cluster
//	@Tags			Tree Cluster
//	@Produce		json
//	@Success		200	{object}	treecluster.TreeClusterResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/cluster/{treecluster_id} [put]
//	@Param			treecluster_id	path	string									true	"Tree Cluster ID"
//	@Param			body			body	treecluster.TreeClusterUpdateRequest	true	"Tree Cluster Update Request"
//	@Param			Authorization	header	string									true	"Insert your access token"	default(Bearer <Add access token here>)
func UpdateTreeCluster(svc service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.JSON(treecluster.TreeClusterResponse{})
	}
}

//	@Summary		Delete tree cluster
//	@Description	Delete tree cluster
//	@Id				delete-tree-cluster
//	@Tags			Tree Cluster
//	@Produce		json
//	@Success		200	{object}	treecluster.TreeClusterResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/cluster/{treecluster_id} [delete]
//	@Param			treecluster_id	path	string	true	"Tree Cluster ID"
//	@Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func DeleteTreeCluster(svc service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.JSON(treecluster.TreeClusterResponse{})
	}
}

//	@Summary		Get all trees in tree cluster
//	@Description	Get all trees in tree cluster
//	@Id				get-all-trees-in-tree-cluster
//	@Tags			Tree Cluster
//	@Produce		json
//	@Success		200	{object}	tree.TreeResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/cluster/{treecluster_id}/trees [get]
//	@Param			treecluster_id	path	string	true	"Tree Cluster ID"
//	@Param			page			query	string	false	"Page"
//	@Param			limit			query	string	false	"Limit"
//	@Param			age				query	string	false	"Age"
//	@Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func GetTreesInTreeCluster(svc service.Service) fiber.Handler {
	// TODO: Change response @Success to tree.TreeListResponse
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.JSON(tree.TreeResponse{})
	}
}

//	@Summary		Add trees to tree cluster
//	@Description	Add trees to tree cluster
//	@Id				add-trees-to-tree-cluster
//	@Tags			Tree Cluster
//	@Produce		json
//	@Success		200	{object}	treecluster.TreeClusterResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/cluster/{treecluster_id}/trees [post]
//	@Param			treecluster_id	path	string									true	"Tree Cluster ID"
//	@Param			body			body	treecluster.TreeClusterAddTreesRequest	true	"Tree Cluster Add Trees Request"
//	@Param			Authorization	header	string									true	"Insert your access token"	default(Bearer <Add access token here>)
func AddTreesToTreeCluster(svc service.Service) fiber.Handler {
  return func(c *fiber.Ctx) error {
    // TODO: Implement
    return c.JSON(treecluster.TreeClusterResponse{})
  }
}

//	@Summary		Remove trees from tree cluster
//	@Description	Remove trees from tree cluster
//	@Id				remove-trees-from-tree-cluster
//	@Tags			Tree Cluster
//	@Produce		json
//	@Success		200	{object}	treecluster.TreeClusterResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/cluster/{treecluster_id}/trees/{tree_id} [delete]
//	@Param			treecluster_id	path	string	true	"Tree Cluster ID"
//	@Param			tree_id			path	string	true	"Tree ID"
//	@Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func RemoveTreesFromTreeCluster(svc service.Service) fiber.Handler {
  return func(c *fiber.Ctx) error {
    // TODO: Implement
    return c.JSON(treecluster.TreeClusterResponse{})
  }
}


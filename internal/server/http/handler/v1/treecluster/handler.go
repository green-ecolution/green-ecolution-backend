package treecluster

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/errorhandler"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

var (
	treeClusterMapper = generated.TreeClusterHTTPMapperImpl{}
	treeMapper        = generated.TreeHTTPMapperImpl{}
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
func GetAllTreeClusters(svc service.TreeClusterService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		domainData, err := svc.GetAll(ctx)
		if err != nil {
			return errorhandler.HandleError(err)
		}

		data := make([]*entities.TreeClusterResponse, len(domainData))
		for i, domain := range domainData {
			data[i] = mapTreeClusterToDto(domain)
		}

		return c.JSON(entities.TreeClusterListResponse{
			Data:       data,
			Pagination: &entities.Pagination{}, // TODO: Handle pagination
		})
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
func GetTreeClusterByID(svc service.TreeClusterService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		id, err := strconv.Atoi(c.Params("treecluster_id"))
		if err != nil {
			return err
		}

		domainData, err := svc.GetByID(ctx, int32(id))

		if err != nil {
			return errorhandler.HandleError(err)
		}

		data := mapTreeClusterToDto(domainData)

		return c.JSON(data)
	}
}

// @Summary		Create tree cluster
// @Description	Create tree cluster
// @Id				create-tree-cluster
// @Tags			Tree Cluster
// @Produce		json
// @Success		201	{object}	entities.TreeClusterResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/cluster [post]
// @Param			body			body	entities.TreeClusterCreateRequest	true	"Tree Cluster Create Request"
// @Param			Authorization	header	string								true	"Insert your access token"	default(Bearer <Add access token here>)
func CreateTreeCluster(svc service.TreeClusterService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		var req entities.TreeClusterCreateRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		domainReq := treeClusterMapper.FromCreateRequest(&req)
		domainData, err := svc.Create(ctx, domainReq)
		if err != nil {
			return errorhandler.HandleError(err)
		}

		data := mapTreeClusterToDto(domainData)
		return c.Status(fiber.StatusCreated).JSON(data)
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
func UpdateTreeCluster(svc service.TreeClusterService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		id, err := strconv.Atoi(c.Params("treecluster_id"))
		if err != nil {
			return errorhandler.HandleError(err)
		}

		var req entities.TreeClusterUpdateRequest
		if err = c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		domainReq := treeClusterMapper.FromUpdateRequest(&req)
		domainData, err := svc.Update(ctx, int32(id), domainReq)
		if err != nil {
			return errorhandler.HandleError(err)
		}

		data := mapTreeClusterToDto(domainData)
		return c.JSON(data)
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
func DeleteTreeCluster(svc service.TreeClusterService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		id, err := strconv.Atoi(c.Params("treecluster_id"))
		if err != nil {
			return err
		}

		err = svc.Delete(ctx, int32(id))
		if err != nil {
			return errorhandler.HandleError(err)
		}

		return c.SendStatus(fiber.StatusNoContent)
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
func GetTreesInTreeCluster(_ service.TreeClusterService) fiber.Handler {
	// TODO: Change response @Success to entities.TreeListResponse
	return func(c *fiber.Ctx) error {
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
func AddTreesToTreeCluster(_ service.TreeClusterService) fiber.Handler {
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
func RemoveTreesFromTreeCluster(_ service.TreeClusterService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

func mapTreeClusterToDto(t *domain.TreeCluster) *entities.TreeClusterResponse {
	dto := treeClusterMapper.FormResponse(t)

	if t.Region != nil {
		dto.Region = &entities.RegionResponse{
			ID:   t.Region.ID,
			Name: t.Region.Name,
		}
	}

	dto.Trees = treeMapper.FromResponseList(t.Trees)

	return dto
}

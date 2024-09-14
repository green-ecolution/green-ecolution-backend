package tree

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/tree"
	_ "github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/errorhandler"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

// @Summary		Get all trees
// @Description	Get all trees
// @Id				get-all-trees
// @Tags			Trees
// @Produce		json
// @Success		200	{object}	tree.TreeListResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/tree [get]
// @Param			page			query	string	false	"Page"
// @Param			limit			query	string	false	"Limit"
// @Param			age				query	string	false	"Age"
// @Param			treecluster_id	query	string	false	"Tree Cluster ID"
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func GetAllTrees(svc service.TreeService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.JSON(tree.TreeResponse{})
	}
}

// @Summary		Get tree by ID
// @Description	Get tree by ID
// @Id				get-trees
// @Tags			Trees
// @Produce		json
// @Success		200	{object}	tree.TreeResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/tree/{tree_id} [get]
// @Param			tree_id			path	string	false	"Tree ID"
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func GetTreeByID(svc service.TreeService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.JSON(tree.TreeResponse{})
	}
}

// @Summary		Create tree
// @Description	Create tree
// @Id				create-tree
// @Tags			Trees
// @Produce		json
// @Success		200	{object}	tree.TreeResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/tree [post]
// @Param			Authorization	header	string					true	"Insert your access token"	default(Bearer <Add access token here>)
// @Param			body			body	tree.TreeCreateRequest	true	"Tree to create"
func CreateTree(svc service.TreeService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.JSON(tree.TreeResponse{})
	}
}

// @Summary		Update tree
// @Description	Update tree
// @Id				update-tree
// @Tags			Trees
// @Produce		json
// @Success		200	{object}	tree.TreeResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/tree/{tree_id} [patch]
// @Param			Authorization	header	string					true	"Insert your access token"	default(Bearer <Add access token here>)
// @Param			tree_id			path	string					false	"Tree ID"
// @Param			body			body	tree.TreeUpdateRequest	true	"Tree to update"
func UpdateTree(svc service.TreeService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.JSON(tree.TreeResponse{})
	}
}

// @Summary		Delete tree
// @Description	Delete tree
// @Id				delete-tree
// @Tags			Trees
// @Produce		json
// @Success		200
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/tree [delete]
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func DeleteTree(svc service.TreeService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.JSON(tree.TreeResponse{})
	}
}

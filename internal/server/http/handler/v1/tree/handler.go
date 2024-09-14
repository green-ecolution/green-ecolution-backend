package tree

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/tree"
	_ "github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/errorhandler"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

//	@Summary		Get all trees
//	@Description	Get all trees
//	@Id				get-all-trees
//	@Tags			Tree
//	@Produce		json
//	@Success		200	{object}	tree.TreeListResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/tree [get]
//	@Param			page			query	string	false	"Page"
//	@Param			limit			query	string	false	"Limit"
//	@Param			age				query	string	false	"Age"
//	@Param			treecluster_id	query	string	false	"Tree Cluster ID"
//	@Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func GetAllTrees(_ service.TreeService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.JSON(tree.TreeResponse{})
	}
}

//	@Summary		Get tree by ID
//	@Description	Get tree by ID
//	@Id				get-trees
//	@Tags			Tree
//	@Produce		json
//	@Success		200	{object}	tree.TreeResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/tree/{tree_id} [get]
//	@Param			tree_id			path	string	false	"Tree ID"
//	@Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func GetTreeByID(_ service.TreeService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.JSON(tree.TreeResponse{})
	}
}

//	@Summary		Create tree
//	@Description	Create tree
//	@Id				create-tree
//	@Tags			Tree
//	@Produce		json
//	@Success		200	{object}	tree.TreeResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/tree [post]
//	@Param			Authorization	header	string					true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			body			body	tree.TreeCreateRequest	true	"Tree to create"
func CreateTree(_ service.TreeService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.JSON(tree.TreeResponse{})
	}
}

//	@Summary		Update tree
//	@Description	Update tree
//	@Id				update-tree
//	@Tags			Tree
//	@Produce		json
//	@Success		200	{object}	tree.TreeResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/tree/{tree_id} [put]
//	@Param			Authorization	header	string					true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			tree_id			path	string					false	"Tree ID"
//	@Param			body			body	tree.TreeUpdateRequest	true	"Tree to update"
func UpdateTree(_ service.TreeService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.JSON(tree.TreeResponse{})
	}
}

//	@Summary		Delete tree
//	@Description	Delete tree
//	@Id				delete-tree
//	@Tags			Tree
//	@Produce		json
//	@Success		200
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/tree [delete]
//	@Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func DeleteTree(_ service.TreeService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.JSON(tree.TreeResponse{})
	}
}

//	@Summary		Get sensors of a tree
//	@Description	Get sensors of a tree
//	@Id				get-tree-sensors
//	@Tags			Tree Sensor
//	@Produce		json
//	@Success		200	{object}	tree.TreeResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/tree/{tree_id}/sensors [get]
//	@Param			tree_id			path	string	false	"Tree ID"
//	@Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func GetTreeSensor(_ service.TreeService) fiber.Handler {
  // TODO: Change @Success to return sensor.SensorResponse
  return func(c *fiber.Ctx) error {
    // TODO: Implement
    return c.JSON(tree.TreeResponse{})
  }
}

//	@Summary		Add sensor to a tree
//	@Description	Add sensor to a tree
//	@Id				add-sensor-to-tree
//	@Tags			Tree Sensor
//	@Produce		json
//	@Success		200	{object}	tree.TreeResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/tree/{tree_id}/sensors [post]
//	@Param			tree_id			path	string						false	"Tree ID"
//	@Param			body			body	tree.TreeAddSensorRequest	true	"Sensor to add"
//	@Param			Authorization	header	string						true	"Insert your access token"	default(Bearer <Add access token here>)
func AddTreeSensor(_ service.TreeService) fiber.Handler {
  return func(c *fiber.Ctx) error {
    // TODO: Implement
    return c.JSON(tree.TreeResponse{})
  }
}

//	@Summary		Remove sensor from a tree
//	@Description	Remove sensor from a tree
//	@Id				remove-sensor-from-tree
//	@Tags			Tree Sensor
//	@Produce		json
//	@Success		200	{object}	tree.TreeResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/tree/{tree_id}/sensors/{sensor_id} [delete]
//	@Param			tree_id			path	string	false	"Tree ID"
//	@Param			sensor_id		path	string	false	"Sensor ID"
//	@Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func RemoveTreeSensor(_ service.TreeService) fiber.Handler {
  return func(c *fiber.Ctx) error {
    // TODO: Implement
    return c.JSON(tree.TreeResponse{})
  }
}

//	@Summary		Get images of a tree
//	@Description	Get images of a tree
//	@Id				get-tree-images
//	@Tags			Tree Images
//	@Produce		json
//	@Success		200	{object}	tree.TreeResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/tree/{tree_id}/images [get]
//	@Param			tree_id			path	string	false	"Tree ID"
//	@Param			page			query	string	false	"Page"
//	@Param			limit			query	string	false	"Limit"
//	@Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func GetTreeImages(_ service.TreeService) fiber.Handler {
  // TODO: Change @Success to return image.ImageResponse
  return func(c *fiber.Ctx) error {
    // TODO: Implement
    return c.JSON(tree.TreeResponse{})
  }
}

//	@Summary		Add images to a tree
//	@Description	Add images to a tree
//	@Id				add-images-to-tree
//	@Tags			Tree Images
//	@Produce		json
//	@Success		200	{object}	tree.TreeResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/tree/{tree_id}/images [post]
//	@Param			tree_id			path	string						false	"Tree ID"
//	@Param			body			body	tree.TreeAddImagesRequest	true	"Images to add"
//	@Param			Authorization	header	string						true	"Insert your access token"	default(Bearer <Add access token here>)
func AddTreeImage(_ service.TreeService) fiber.Handler {
  return func(c *fiber.Ctx) error {
    // TODO: Implement
    return c.JSON(tree.TreeResponse{})
  }
}

//	@Summary		Remove image from a tree
//	@Description	Remove image from a tree
//	@Id				remove-image-from-tree
//	@Tags			Tree Images
//	@Produce		json
//	@Success		200	{object}	tree.TreeResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/tree/{tree_id}/images/{image_id} [delete]
//	@Param			tree_id			path	string	false	"Tree ID"
//	@Param			image_id		path	string	false	"Image ID"
//	@Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func RemoveTreeImage(_ service.TreeService) fiber.Handler {
  return func(c *fiber.Ctx) error {
    // TODO: Implement
    return c.JSON(tree.TreeResponse{})
  }
}


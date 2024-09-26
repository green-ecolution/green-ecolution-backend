package tree

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/errorhandler"
	_ "github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/errorhandler"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

var (
  treeMapper = generated.TreeHTTPMapperImpl{}
  sensorMapper = generated.SensorHTTPMapperImpl{}
)

// @Summary		Get all trees
// @Description	Get all trees
// @Id				get-all-trees
// @Tags			Tree
// @Produce		json
// @Success		200	{object}	entities.TreeListResponse
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
    ctx := c.Context()
    domainData, err := svc.GetAll(ctx)
    if err != nil {
      return errorhandler.HandleError(err)
    }

    data := make([]*entities.TreeResponse, len(domainData))
    for i, domain := range domainData {
      data[i] = treeMapper.FromResponse(domain)
      data[i].Sensor = sensorMapper.FromResponse(domain.Sensor)
    }

    return c.JSON(entities.TreeListResponse{
      Data:       data,
      Pagination: entities.Pagination{}, // TODO: Handle pagination
    })
	}
}

// @Summary		Get tree by ID
// @Description	Get tree by ID
// @Id				get-trees
// @Tags			Tree
// @Produce		json
// @Success		200	{object}	entities.TreeResponse
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
    ctx := c.Context()
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
      return err
    }

    domainData, err := svc.GetByID(ctx, id)
    if err != nil {
      return errorhandler.HandleError(err)
    }

    data := *treeMapper.FromResponse(domainData)
    data.Sensor = sensorMapper.FromResponse(domainData.Sensor)

    return c.JSON(data)
	}
}

// @Summary		Create tree
// @Description	Create tree
// @Id				create-tree
// @Tags			Tree
// @Produce		json
// @Success		200	{object}	entities.TreeResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/tree [post]
// @Param			Authorization	header	string						true	"Insert your access token"	default(Bearer <Add access token here>)
// @Param			body			body	entities.TreeCreateRequest	true	"Tree to create"
func CreateTree(_ service.TreeService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

// @Summary		Update tree
// @Description	Update tree
// @Id				update-tree
// @Tags			Tree
// @Produce		json
// @Success		200	{object}	entities.TreeResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/tree/{tree_id} [put]
// @Param			Authorization	header	string						true	"Insert your access token"	default(Bearer <Add access token here>)
// @Param			tree_id			path	string						false	"Tree ID"
// @Param			body			body	entities.TreeUpdateRequest	true	"Tree to update"
func UpdateTree(_ service.TreeService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

// @Summary		Delete tree
// @Description	Delete tree
// @Id				delete-tree
// @Tags			Tree
// @Produce		json
// @Success		200
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/tree [delete]
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func DeleteTree(_ service.TreeService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

// @Summary		Get sensor of a tree
// @Description	Get sensor of a tree
// @Id				get-tree-sensor
// @Tags			Tree Sensor
// @Produce		json
// @Success		200	{object}	entities.TreeResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/tree/{tree_id}/sensor [get]
// @Param			tree_id			path	string	false	"Tree ID"
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func GetTreeSensor(_ service.TreeService) fiber.Handler {
	// TODO: Change @Success to return sensor.SensorResponse
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

// @Summary		Add sensor to a tree
// @Description	Add sensor to a tree
// @Id				add-sensor-to-tree
// @Tags			Tree Sensor
// @Produce		json
// @Success		200	{object}	entities.TreeResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/tree/{tree_id}/sensor [post]
// @Param			tree_id			path	string							false	"Tree ID"
// @Param			body			body	entities.TreeAddSensorRequest	true	"Sensor to add"
// @Param			Authorization	header	string							true	"Insert your access token"	default(Bearer <Add access token here>)
func AddTreeSensor(_ service.TreeService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

// @Summary		Remove sensor from a tree
// @Description	Remove sensor from a tree
// @Id				remove-sensor-from-tree
// @Tags			Tree Sensor
// @Produce		json
// @Success		200	{object}	entities.TreeResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/tree/{tree_id}/sensor/{sensor_id} [delete]
// @Param			tree_id			path	string	false	"Tree ID"
// @Param			sensor_id		path	string	false	"Sensor ID"
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func RemoveTreeSensor(_ service.TreeService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

// @Summary		Get images of a tree
// @Description	Get images of a tree
// @Id				get-tree-images
// @Tags			Tree Images
// @Produce		json
// @Success		200	{object}	entities.TreeResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/tree/{tree_id}/images [get]
// @Param			tree_id			path	string	false	"Tree ID"
// @Param			page			query	string	false	"Page"
// @Param			limit			query	string	false	"Limit"
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func GetTreeImages(_ service.TreeService) fiber.Handler {
	// TODO: Change @Success to return image.ImageResponse
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

// @Summary		Add images to a tree
// @Description	Add images to a tree
// @Id				add-images-to-tree
// @Tags			Tree Images
// @Produce		json
// @Success		200	{object}	entities.TreeResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/tree/{tree_id}/images [post]
// @Param			tree_id			path	string							false	"Tree ID"
// @Param			body			body	entities.TreeAddImagesRequest	true	"Images to add"
// @Param			Authorization	header	string							true	"Insert your access token"	default(Bearer <Add access token here>)
func AddTreeImage(_ service.TreeService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

// @Summary		Remove image from a tree
// @Description	Remove image from a tree
// @Id				remove-image-from-tree
// @Tags			Tree Images
// @Produce		json
// @Success		200	{object}	entities.TreeResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/tree/{tree_id}/images/{image_id} [delete]
// @Param			tree_id			path	string	false	"Tree ID"
// @Param			image_id		path	string	false	"Image ID"
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func RemoveTreeImage(_ service.TreeService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

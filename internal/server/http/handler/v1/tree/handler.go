package tree

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
	treeMapper   = generated.TreeHTTPMapperImpl{}
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
// @Security		Keycloak
func GetAllTrees(svc service.TreeService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		domainData, err := svc.GetAll(ctx)
		if err != nil {
			return errorhandler.HandleError(err)
		}

		data := make([]*entities.TreeResponse, len(domainData))
		for i, domain := range domainData {
			data[i] = mapTreeToDto(domain)
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
// @Param			tree_id	path	string	false	"Tree ID"
// @Security		Keycloak
func GetTreeByID(svc service.TreeService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			err = service.NewError(service.BadRequest, "invalid ID format")
			return errorhandler.HandleError(err)
		}

		domainData, err := svc.GetByID(ctx, int32(id))
		if err != nil {
			return errorhandler.HandleError(err)
		}

		data := mapTreeToDto(domainData)

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
// @Security		Keycloak
// @Param			body	body	entities.TreeCreateRequest	true	"Tree to create"
func CreateTree(svc service.TreeService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		var req entities.TreeCreateRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		domainReq := treeMapper.FromCreateRequest(&req)
		domainData, err := svc.Create(ctx, domainReq)
		if err != nil {
			return errorhandler.HandleError(err)
		}

		data := mapTreeToDto(domainData)
		return c.Status(fiber.StatusCreated).JSON(data)
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
// @Security		Keycloak
// @Param			tree_id	path	string						false	"Tree ID"
// @Param			body	body	entities.TreeUpdateRequest	true	"Tree to update"
func UpdateTree(svc service.TreeService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			err = service.NewError(service.BadRequest, "invalid ID format")
			return errorhandler.HandleError(err)
		}
		var req entities.TreeUpdateRequest
		if err = c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		domainReq := treeMapper.FromUpdateRequest(&req)
		domainData, err := svc.Update(ctx, int32(id), domainReq)
		if err != nil {
			return errorhandler.HandleError(err)
		}
		data := mapTreeToDto(domainData)
		return c.JSON(data)
	}
}

// @Summary		Delete tree
// @Description	Delete tree
// @Id				delete-tree
// @Tags			Tree
// @Produce		json
// @Success		204
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/tree/{tree_id} [delete]
// @Param			tree_id	path	string	false	"Tree ID"
// @Security		Keycloak
func DeleteTree(svc service.TreeService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			err = service.NewError(service.BadRequest, "invalid ID format")
			return errorhandler.HandleError(err)
		}
		err = svc.Delete(ctx, int32(id))
		if err != nil {
			return errorhandler.HandleError(err)
		}
		return c.SendStatus(fiber.StatusNoContent)
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
// @Param			tree_id	path	string	false	"Tree ID"
// @Security		Keycloak
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
// @Param			tree_id	path	string							false	"Tree ID"
// @Param			body	body	entities.TreeAddSensorRequest	true	"Sensor to add"
// @Security		Keycloak
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
// @Param			tree_id		path	string	false	"Tree ID"
// @Param			sensor_id	path	string	false	"Sensor ID"
// @Security		Keycloak
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
// @Param			tree_id	path	string	false	"Tree ID"
// @Param			page	query	string	false	"Page"
// @Param			limit	query	string	false	"Limit"
// @Security		Keycloak
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
// @Param			tree_id	path	string							false	"Tree ID"
// @Param			body	body	entities.TreeAddImagesRequest	true	"Images to add"
// @Security		Keycloak
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
// @Param			tree_id		path	string	false	"Tree ID"
// @Param			image_id	path	string	false	"Image ID"
// @Security		Keycloak
func RemoveTreeImage(_ service.TreeService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

func mapTreeToDto(t *domain.Tree) *entities.TreeResponse {
	dto := treeMapper.FromResponse(t)
	dto.Sensor = sensorMapper.FromResponse(t.Sensor)

	return dto
}

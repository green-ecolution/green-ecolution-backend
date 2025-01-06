package routing

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/errorhandler"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

// @Summary		Generate preview route
// @Description	Generate preview route
// @Tags			Route
// @Produce		json
// @Accept    json
// @Success		200		{object}	entities.GeoJSON
// @Failure		400		{object}	HTTPError
// @Failure		500		{object}	HTTPError
// @Param			body body	entities.RouteRequest	true	"Route Request"
// @Router			/v1/route/preview [post]
// @Security		Keycloak
func CreatePreviewRoute(svc service.RoutingService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		var req entities.RouteRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		domainGeo, err := svc.PreviewRoute(ctx, req.VehicleID, req.ClusterIDs)
		if err != nil {
			return errorhandler.HandleError(err)
		}

		return c.JSON(entities.GeoJSON{
			Type:     entities.GeoJSONType(domainGeo.Type),
			Bbox:     domainGeo.Bbox,
			Features: domainGeo.Features,
		})
	}
}

// @Summary		Generate route
// @Description	Generate route
// @Tags			Route
// @Produce		json
// @Accept    json
// @Success		200		{object}	entities.GeoJSON "TODO: Change to real entity"
// @Failure		400		{object}	HTTPError
// @Failure		500		{object}	HTTPError
// @Param			body body	entities.RouteRequest	true	"Route Request"
// @Router			/v1/route [post]
// @Security		Keycloak
func CreateRoute(_ service.RoutingService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

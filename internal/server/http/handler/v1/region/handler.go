package region

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	_ "github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/errorhandler"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

// @Summary		Get all regions
// @Description	Get all regions
// @Tags			Region
// @Produce		json
// @Success		200		{object}	entities.RegionListResponse
// @Failure		400		{object}	HTTPError
// @Failure		500		{object}	HTTPError
// @Param			page	query		string	false	"Page"
// @Param			limit	query		string	false	"Limit"
// @Router			/v1/region [get]
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func GetAllRegions(svc service.RegionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
    r, err := svc.GetAll(c.Context())
    if err != nil {
			return errorhandler.HandleError(err)
    }

    return c.JSON(r)
	}
}

// @Summary		Get a region by ID
// @Description	Get a region by ID
// @Tags			Region
// @Produce		json
// @Success		200	{object}	entities.RegionResponse
// @Failure		400	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Param			id	path		string	true	"Region ID"
// @Router			/v1/region/{id} [get]
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func GetRegionByID(svc service.RegionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
      return errorhandler.HandleError(err)
    }

    r, err := svc.GetByID(c.Context(), int32(id))
    if err != nil {
      return errorhandler.HandleError(err)
    }

    return c.JSON(r)
	}
}

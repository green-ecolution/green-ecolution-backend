package vehicle

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/errorhandler"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

var (
	vehicleMapper = generated.VehicleHTTPMapperImpl{}
)

// @Summary		Get all vehicles
// @Description	Get all vehicles
// @Id				get-all-vehicles
// @Tags			Vehicle
// @Produce		json
// @Success		200	{object}	entities.VehicleListResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/vehicle [get]
// @Param			page	query	string	false	"Page"
// @Param			limit	query	string	false	"Limit"
// @Param			status	query	string	false	"Status"
// @Security		Keycloak
func GetAllVehicles(svc service.VehicleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		domainData, err := svc.GetAll(ctx)
		if err != nil {
			return errorhandler.HandleError(err)
		}

		data := make([]*entities.VehicleResponse, len(domainData))
		for i, domain := range domainData {
			data[i] = mapVehicleToDto(domain)
		}

		return c.JSON(entities.VehicleListResponse{
			Data:       data,
			Pagination: &entities.Pagination{}, // TODO: Handle pagination
		})
	}
}

// @Summary		Get vehicle by ID
// @Description	Get vehicle by ID
// @Id				get-vehicle-by-id
// @Tags			Vehicle
// @Produce		json
// @Success		200	{object}	entities.VehicleResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/vehicle/{id} [get]
// @Param			id	path	string	true	"Vehicle ID"
// @Security		Keycloak
func GetVehicleByID(svc service.VehicleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			err := service.NewError(service.BadRequest, "invalid ID format")
			return errorhandler.HandleError(err)
		}

		domainData, err := svc.GetByID(ctx, int32(id))

		if err != nil {
			return errorhandler.HandleError(err)
		}

		data := mapVehicleToDto(domainData)
		return c.JSON(data)
	}
}

// @Summary		Get vehicle by plate
// @Description	Get vehicle by plate
// @Id				get-vehicle-by-plate
// @Tags			Vehicle
// @Produce		json
// @Success		200	{object}	entities.VehicleResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/vehicle/plate/{plate} [get]
// @Param			plate	path	string	true	"Vehicle plate number"
// @Security		Keycloak
func GetVehicleByPlate(svc service.VehicleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		plate := c.Params("plate")
		if plate == "" {
			fmt.Print(plate);
			err := service.NewError(service.BadRequest, "invalid Plate format")
			return errorhandler.HandleError(err)
		}

		domainData, err := svc.GetByPlate(ctx, plate)
		if err != nil {
			return errorhandler.HandleError(err)
		}

		data := mapVehicleToDto(domainData)
		return c.JSON(data)
	}
}

func mapVehicleToDto(v *domain.Vehicle) *entities.VehicleResponse {
	dto := vehicleMapper.FormResponse(v)

	return dto
}

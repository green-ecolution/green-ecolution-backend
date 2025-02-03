package vehicle

import (
	"strconv"
	"strings"

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

//	@Summary		Get all vehicles
//	@Description	Get all vehicles
//	@Id				get-all-vehicles
//	@Tags			Vehicle
//	@Produce		json
//	@Success		200	{object}	entities.VehicleListResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/vehicle [get]
//	@Param			page		query	string	false	"Page"
//	@Param			limit		query	string	false	"Limit"
//	@Param			status		query	string	false	"Status"
//	@Param			type		query	string	false	"Vehicle Type"
//	@Param			provider	query	string	false	"Provider"
//	@Security		Keycloak
func GetAllVehicles(svc service.VehicleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		var domainData []*domain.Vehicle
		var err error

		vehicleTypeStr := c.Query("type")
		provider := c.Query("provider")
		if vehicleTypeStr != "" {
			vehicleType := domain.ParseVehicleType(vehicleTypeStr)
			if vehicleType == domain.VehicleTypeUnknown {
				return errorhandler.HandleError(service.NewError(service.BadRequest, "invalid vehicle type"))
			}
			domainData, err = svc.GetAllByType(ctx, vehicleType)
		} else {
			domainData, err = svc.GetAll(ctx, provider)
		} // TODO: move to service

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

//	@Summary		Get vehicle by ID
//	@Description	Get vehicle by ID
//	@Id				get-vehicle-by-id
//	@Tags			Vehicle
//	@Produce		json
//	@Success		200	{object}	entities.VehicleResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/vehicle/{id} [get]
//	@Param			id	path	string	true	"Vehicle ID"
//	@Security		Keycloak
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

//	@Summary		Get vehicle by plate
//	@Description	Get vehicle by plate
//	@Id				get-vehicle-by-plate
//	@Tags			Vehicle
//	@Produce		json
//	@Success		200	{object}	entities.VehicleResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/vehicle/plate/{plate} [get]
//	@Param			plate	path	string	true	"Vehicle plate number"
//	@Security		Keycloak
func GetVehicleByPlate(svc service.VehicleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		plate := strings.Clone(c.Params("plate"))
		if plate == "" {
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

//	@Summary		Create vehicle
//	@Description	Create vehicle
//	@Id				create-vehicle
//	@Tags			Vehicle
//	@Produce		json
//	@Success		201	{object}	entities.VehicleResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/vehicle [post]
//	@Param			body	body	entities.VehicleCreateRequest	true	"Vehicle Create Request"
//	@Security		Keycloak
func CreateVehicle(svc service.VehicleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		var req entities.VehicleCreateRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		domainReq := vehicleMapper.FromCreateRequest(&req)
		domainData, err := svc.Create(ctx, domainReq)
		if err != nil {
			return errorhandler.HandleError(err)
		}

		data := mapVehicleToDto(domainData)
		return c.Status(fiber.StatusCreated).JSON(data)
	}
}

//	@Summary		Update vehicle
//	@Description	Update vehicle
//	@Id				update-vehicle
//	@Tags			Vehicle
//	@Produce		json
//	@Success		200	{object}	entities.VehicleResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/vehicle/{id} [put]
//	@Param			id		path	string							true	"Vehicle ID"
//	@Param			body	body	entities.VehicleUpdateRequest	true	"Vehicle Update Request"
//	@Security		Keycloak
func UpdateVehicle(svc service.VehicleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			err := service.NewError(service.BadRequest, "invalid ID format")
			return errorhandler.HandleError(err)
		}

		var req entities.VehicleUpdateRequest
		if err = c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		domainReq := vehicleMapper.FromUpdateRequest(&req)
		domainData, err := svc.Update(ctx, int32(id), domainReq)
		if err != nil {
			return errorhandler.HandleError(err)
		}

		data := mapVehicleToDto(domainData)
		return c.JSON(data)
	}
}

//	@Summary		Delete vehicle
//	@Description	Delete vehicle
//	@Id				delete-vehicle
//	@Tags			Vehicle
//	@Produce		json
//	@Success		204
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/vehicle/{id} [delete]
//	@Param			id	path	string	true	"Vehicle ID"
//	@Security		Keycloak
func DeleteVehicle(svc service.VehicleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			err := service.NewError(service.BadRequest, "invalid ID format")
			return errorhandler.HandleError(err)
		}

		err = svc.Delete(ctx, int32(id))
		if err != nil {
			return errorhandler.HandleError(err)
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

func mapVehicleToDto(v *domain.Vehicle) *entities.VehicleResponse {
	dto := vehicleMapper.FromResponse(v)

	return dto
}

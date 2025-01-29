package sensor

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/errorhandler"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

var (
	sensorMapper = generated.SensorHTTPMapperImpl{}
)

//	@Summary		Get all sensors
//	@Description	Get all sensors
//	@Id				get-all-sensors
//	@Tags			Sensor
//	@Produce		json
//	@Success		200	{object}	entities.SensorListResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/sensor [get]
//	@Param			status		query	string	false	"Sensor Status"
//	@Param			page		query	string	false	"Page"
//	@Param			limit		query	string	false	"Limit"
//	@Param			provider	query	string	false	"Provider"
//	@Security		Keycloak
func GetAllSensors(svc service.SensorService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		provider := c.Query("provider")
		domainData, err := svc.GetAll(ctx, provider)
		if err != nil {
			return errorhandler.HandleError(err)
		}

		data := make([]*entities.SensorResponse, len(domainData))
		for i, domain := range domainData {
			data[i] = mapToDto(domain)
		}

		return c.JSON(entities.SensorListResponse{
			Data:       data,
			Pagination: entities.Pagination{}, // TODO: Handle Pagination
		})
	}
}

//	@Summary		Get sensor by ID
//	@Description	Get sensor by ID
//	@Id				get-sensor-by-id
//	@Tags			Sensor
//	@Produce		json
//	@Success		200	{object}	entities.SensorResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/sensor/{sensor_id} [get]
//	@Param			sensor_id	path	string	true	"Sensor ID"
//	@Security		Keycloak
func GetSensorByID(svc service.SensorService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		id := strings.Clone(c.Params("id"))
		if id == "" {
			err := service.NewError(service.BadRequest, "invalid ID format")
			return errorhandler.HandleError(err)
		}

		domainData, err := svc.GetByID(ctx, id)

		if err != nil {
			return errorhandler.HandleError(err)
		}

		data := mapToDto(domainData)

		return c.JSON(data)
	}
}

//	@Summary		Delete sensor
//	@Description	Delete sensor
//	@Id				delete-sensor
//	@Tags			Sensor
//	@Produce		json
//	@Success		204
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/sensor/{sensor_id} [delete]
//	@Param			sensor_id	path	string	true	"Sensor ID"
//	@Security		Keycloak
func DeleteSensor(svc service.SensorService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		id := strings.Clone(c.Params("id"))
		if id == "" {
			err := service.NewError(service.BadRequest, "invalid ID format")
			return errorhandler.HandleError(err)
		}

		err := svc.Delete(ctx, id)

		if err != nil {
			return errorhandler.HandleError(err)
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

func mapToDto(t *domain.Sensor) *entities.SensorResponse {
	dto := sensorMapper.FromResponse(t)
	return dto
}

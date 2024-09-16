package sensor

import (
	"github.com/gofiber/fiber/v2"
	_ "github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

// @Summary		Get all sensors
// @Description	Get all sensors
// @Id				get-all-sensors
// @Tags			Sensor
// @Produce		json
// @Success		200	{object}	[]entities.SensorListResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/sensor [get]
// @Param			status			query	string	false	"Sensor Status"
// @Param			sensor_id		path	string	true	"Sensor ID"
// @Param			page			query	string	false	"Page"
// @Param			limit			query	string	false	"Limit"
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func GetAllSensor(_ service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

// @Summary		Get sensor by ID
// @Description	Get sensor by ID
// @Id				get-sensor-by-id
// @Tags			Sensor
// @Produce		json
// @Success		200	{object}	entities.SensorResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/sensor/{sensor_id} [get]
// @Param			sensor_id		path	string	true	"Sensor ID"
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func GetSensorByID(_ service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

// @Summary		Get sensor data by ID
// @Description	Get sensor data by ID
// @Id				get-sensor-data-by-id
// @Tags			Sensor
// @Produce		json
// @Success		200	{object}	entities.SensorDataListResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/sensor/{sensor_id}/data [get]
// @Param			sensor_id		path	string	true	"Sensor ID"
// @Param			page			query	string	false	"Page"
// @Param			limit			query	string	false	"Limit"
// @Param			start_time		query	string	false	"Start time"
// @Param			end_time		query	string	false	"End time"
// @Param			treecluster_id	query	string	false	"TreeCluster ID"
// @Param			Authorization	header	string	false	"Insert your access token"	default(Bearer <Add access token here>)
func GetSensorDataByID(_ service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

// @Summary		Create sensor
// @Description	Create sensor
// @Id				create-sensor
// @Tags			Sensor
// @Produce		json
// @Success		200	{object}	entities.SensorResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/sensor/ [post]
// @Param			Authorization	header	string							false	"Insert your access token"	default(Bearer <Add access token here>)
// @Param			body			body	entities.SensorCreateRequest	true	"Sensor to create"
func CreateSensor(_ service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

// @Summary		Update sensor
// @Description	Update sensor
// @Id				update-sensor
// @Tags			Sensor
// @Produce		json
// @Success		200	{object}	entities.SensorResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/sensor/{sensor_id} [put]
// @Param			sensor_id		path	string							true	"Sensor ID"
// @Param			Authorization	header	string							false	"Insert your access token"	default(Bearer <Add access token here>)
// @Param			body			body	entities.SensorUpdateRequest	true	"Sensor information to update"
func UpdateSensor(_ service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

// @Summary		Delete sensor
// @Description	Delete sensor
// @Id				delete-sensor
// @Tags			Sensor
// @Produce		json
// @Success		200
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/sensor/{sensor_id} [delete]
// @Param			sensor_id		path	string	true	"Sensor ID"
// @Param			Authorization	header	string	false	"Insert your access token"	default(Bearer <Add access token here>)
func DeleteSensor(_ service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement
		return c.SendString("Not implemented")
	}
}

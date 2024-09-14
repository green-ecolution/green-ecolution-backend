package sensor

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/sensor"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

// @Summary		Get all sensors
// @Description	Get all sensors
// @Id				get-all-sensors
// @Tags			Sensor
// @Produce		json
// @Success		200	{object}	[]sensor.SensorResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/sensor [get]
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func GetAllSensor(_ service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement GetAll
		return c.JSON([]sensor.SensorResponse{})
	}
}

// @Summary		Get sensor by ID
// @Description	Get sensor by ID
// @Id				get-sensor-by-id
// @Tags			Sensor
// @Produce		json
// @Success		200	{object}	sensor.SensorResponse
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
		// TODO: Implement GetByID
		return c.JSON(sensor.SensorResponse{})
	}
}

// @Summary		Get sensor data by ID
// @Description	Get sensor data by ID
// @Id				get-sensor-data-by-id
// @Tags			Sensor
// @Produce		json
// @Success		200	{object}	sensor.SensorDataResponse
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
// @Param			Authorization	header	string	false	"Insert your access token"	default(Bearer <Add access token here>)
func GetSensorDataByID(_ service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement GetByID
		return c.JSON(sensor.SensorDataResponse{})
	}
}

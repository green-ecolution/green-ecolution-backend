package sensor_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	serverEntities "github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/sensor"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	serviceMock "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllSensors(t *testing.T) {
	t.Run("should return all sensors successfully with full MqttPayload", func(t *testing.T) {
		mockSensorService := serviceMock.NewMockSensorService(t)
		app := fiber.New()
		handler := sensor.GetAllSensors(mockSensorService)

		mockSensorService.EXPECT().GetAll(
			mock.Anything,
			"",
		).Return(TestSensorList, nil)

		app.Get("/v1/sensor", handler)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/sensor", nil)
		resp, err := app.Test(req, -1)
		assert.Nil(t, err)
		defer resp.Body.Close()

		// then
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response serverEntities.SensorListResponse
		err = utils.ParseJSONResponse(resp, &response)
		assert.NoError(t, err)

		// Assert response matches test data
		assert.Len(t, response.Data, len(TestSensorList))
		assert.Equal(t, TestSensorList[0].ID, response.Data[0].ID)

		mockSensorService.AssertExpectations(t)
	})

	t.Run("should return empty sensor list when no sensors found", func(t *testing.T) {
		mockSensorService := serviceMock.NewMockSensorService(t)
		app := fiber.New()
		handler := sensor.GetAllSensors(mockSensorService)

		mockSensorService.EXPECT().GetAll(
			mock.Anything,
			"",
		).Return([]*entities.Sensor{}, nil)

		app.Get("/v1/sensor", handler)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/sensor", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response serverEntities.SensorListResponse
		err = utils.ParseJSONResponse(resp, &response)
		assert.NoError(t, err)
		assert.Len(t, response.Data, 0)

		mockSensorService.AssertExpectations(t)
	})

	t.Run("should return 500 when service returns an error", func(t *testing.T) {
		mockSensorService := serviceMock.NewMockSensorService(t)
		app := fiber.New()
		handler := sensor.GetAllSensors(mockSensorService)

		mockSensorService.EXPECT().GetAll(
			mock.Anything,
			"",
		).Return(nil, errors.New("service error"))

		app.Get("/v1/sensor", handler)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/sensor", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockSensorService.AssertExpectations(t)
	})
}

func TestGetSensorById(t *testing.T) {
	t.Run("should return sensor by id successfully", func(t *testing.T) {
		mockSensorService := serviceMock.NewMockSensorService(t)
		app := fiber.New()
		handler := sensor.GetSensorByID(mockSensorService)

		mockSensorService.EXPECT().GetByID(
			mock.Anything,
			"sensor-1",
		).Return(TestSensor, nil)

		app.Get("/v1/sensor/:id", handler)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/sensor/sensor-1", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response serverEntities.SensorResponse
		err = utils.ParseJSONResponse(resp, &response)
		assert.NoError(t, err)

		assert.Equal(t, response.ID, TestSensor.ID)
		assert.WithinDuration(t, response.CreatedAt, TestSensor.CreatedAt, time.Second)
		assert.WithinDuration(t, response.UpdatedAt, TestSensor.UpdatedAt, time.Second)
		assert.Equal(t, entities.SensorStatus(response.Status), TestSensor.Status)

		// TODO: compare data

		mockSensorService.AssertExpectations(t)
	})

	t.Run("should return 404 Not Found if sensor does not exist", func(t *testing.T) {
		mockSensorService := serviceMock.NewMockSensorService(t)
		app := fiber.New()
		handler := sensor.GetSensorByID(mockSensorService)

		mockSensorService.EXPECT().GetByID(
			mock.Anything,
			"sensor-999",
		).Return(nil, service.NewError(service.NotFound, "not found"))

		app.Get("/v1/sensor/:id", handler)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/sensor/sensor-999", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		mockSensorService.AssertExpectations(t)
	})

	t.Run("should return 500 Internal Server Error for service failure", func(t *testing.T) {
		mockSensorService := serviceMock.NewMockSensorService(t)
		app := fiber.New()
		handler := sensor.GetSensorByID(mockSensorService)

		mockSensorService.EXPECT().GetByID(
			mock.Anything,
			"sensor-1",
		).Return(nil, fiber.NewError(fiber.StatusInternalServerError, "service error"))

		app.Get("/v1/sensor/:id", handler)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/sensor/sensor-1", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockSensorService.AssertExpectations(t)
	})
}

func TestDeleteSensor(t *testing.T) {
	t.Run("should delete sensor successfully", func(t *testing.T) {
		mockSensorService := serviceMock.NewMockSensorService(t)
		app := fiber.New()
		handler := sensor.DeleteSensor(mockSensorService)

		mockSensorService.EXPECT().Delete(
			mock.Anything,
			"sensor-1",
		).Return(nil)

		app.Delete("/v1/sensor/:id", handler)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, "/v1/sensor/sensor-1", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		mockSensorService.AssertExpectations(t)
	})

	t.Run("should return 404 for non-existing sensor", func(t *testing.T) {
		mockSensorService := serviceMock.NewMockSensorService(t)
		app := fiber.New()
		handler := sensor.DeleteSensor(mockSensorService)

		mockSensorService.EXPECT().Delete(
			mock.Anything,
			"sensor-999",
		).Return(service.NewError(service.NotFound, "not found"))

		app.Delete("/v1/sensor/:id", handler)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, "/v1/sensor/sensor-999", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()
		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

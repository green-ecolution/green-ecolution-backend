package sensor

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	serverEntities "github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	serviceMock "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllSensors(t *testing.T) {
	t.Run("should return all sensors successfully with full MqttPayload", func(t *testing.T) {
		mockSensorService := serviceMock.NewMockSensorService(t)
		app := fiber.New()
		handler := GetAllSensors(mockSensorService)

		mockSensorService.EXPECT().GetAll(
			mock.Anything,
		).Return(TestSensorList, nil)

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

		assert.Len(t, response.Data, len(TestSensorList))
		sensorData := response.Data[0]
		assert.Equal(t, TestSensorList[0].ID, sensorData.ID)

		mockSensorService.AssertExpectations(t)
	})

	t.Run("should return empty sensor list when no sensors found", func(t *testing.T) {
		mockSensorService := serviceMock.NewMockSensorService(t)
		app := fiber.New()
		handler := GetAllSensors(mockSensorService)

		mockSensorService.EXPECT().GetAll(
			mock.Anything,
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
		handler := GetAllSensors(mockSensorService)

		mockSensorService.EXPECT().GetAll(
			mock.Anything,
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
		handler := GetSensorByID(mockSensorService)

		mockSensorService.EXPECT().GetByID(
			mock.Anything,
			int32(1),
		).Return(TestSensor, nil)

		app.Get("/v1/sensor/:id", handler)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/sensor/1", nil)
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
		assert.Equal(t, entities.SensorStatus(response.Status), entities.SensorStatus(TestSensor.Status))
		// TODO: compare data
		mockSensorService.AssertExpectations(t)
	})
}
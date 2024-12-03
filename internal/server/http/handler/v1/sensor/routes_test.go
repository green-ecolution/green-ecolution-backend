package sensor

import (
	"context"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	serviceMock "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterRoutes(t *testing.T) {
	t.Run("/v1/sensor", func(t *testing.T) {
		t.Run("should call GET handler", func(t *testing.T) {
			mockSensorService := serviceMock.NewMockSensorService(t)
			app := fiber.New()
			RegisterRoutes(app, mockSensorService)

			mockSensorService.EXPECT().GetAll(
				mock.Anything,
			).Return(TestSensorList, nil)

			// when
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)

			// then
			resp, err := app.Test(req)
			defer resp.Body.Close()
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	})

	t.Run("/v1/sensor/:id", func(t *testing.T) {
		t.Run("should call GET handler", func(t *testing.T) {
			mockSensorService := serviceMock.NewMockSensorService(t)
			app := fiber.New()
			RegisterRoutes(app, mockSensorService)

			mockSensorService.EXPECT().GetByID(
				mock.Anything,
				"sensor-1",
			).Return(TestSensor, nil)

			// when
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/sensor-1", nil)

			// then
			resp, err := app.Test(req)
			defer resp.Body.Close()
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})

		t.Run("should call DELETE handler", func(t *testing.T) {
			mockSensorService := serviceMock.NewMockSensorService(t)
			app := fiber.New()
			RegisterRoutes(app, mockSensorService)

			mockSensorService.EXPECT().Delete(
				mock.Anything,
				"sensor-1",
			).Return(nil)

			// when
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, "/sensor-1", nil)

			// then
			resp, err := app.Test(req)
			defer resp.Body.Close()
			assert.NoError(t, err)
			assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		})
	})

	t.Run("/v1/sensor/:id/data", func(t *testing.T) {
		t.Run("should call GET handler", func(t *testing.T) {
			mockSensorService := serviceMock.NewMockSensorService(t)
			app := fiber.New()
			RegisterRoutes(app, mockSensorService)

			// when
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/sensor-1/data", nil)

			// then
			resp, err := app.Test(req)
			defer resp.Body.Close()
			assert.NoError(t, err)
			assert.Equal(t, http.StatusNotImplemented, resp.StatusCode)
		})
	})
}

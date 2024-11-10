package sensor_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/sensor"
	serviceMock "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterRoutes(t *testing.T) {
	t.Run("/v1/sensor", func(t *testing.T) {
		t.Run("should call GET handler", func(t *testing.T) {
			mockSensorService := serviceMock.NewMockSensorService(t)
			app := sensor.RegisterRoutes(mockSensorService)

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

		t.Run("should call POST handler", func(t *testing.T) {
			mockSensorService := serviceMock.NewMockSensorService(t)
			app := sensor.RegisterRoutes(mockSensorService)

			// when
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/", nil)

			// then
			resp, err := app.Test(req)
			defer resp.Body.Close()
			assert.NoError(t, err)
			assert.Equal(t, http.StatusNotImplemented, resp.StatusCode)
		})
	})

	t.Run("/v1/sensor/:id", func(t *testing.T) {
		t.Run("should call GET handler", func(t *testing.T) {
			mockSensorService := serviceMock.NewMockSensorService(t)
			app := sensor.RegisterRoutes(mockSensorService)

			// when
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/1", nil)

			// then
			resp, err := app.Test(req)
			defer resp.Body.Close()
			assert.NoError(t, err)
			assert.Equal(t, http.StatusNotImplemented, resp.StatusCode)
		})

		t.Run("should call PUT handler", func(t *testing.T) {
			mockSensorService := serviceMock.NewMockSensorService(t)
			app := sensor.RegisterRoutes(mockSensorService)

			// when
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, "/1", nil)

			// then
			resp, err := app.Test(req)
			defer resp.Body.Close()
			assert.NoError(t, err)
			assert.Equal(t, http.StatusNotImplemented, resp.StatusCode)
		})

		t.Run("should call DELETE handler", func(t *testing.T) {
			mockSensorService := serviceMock.NewMockSensorService(t)
			app := sensor.RegisterRoutes(mockSensorService)

			// when
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, "/1", nil)

			// then
			resp, err := app.Test(req)
			defer resp.Body.Close()
			assert.NoError(t, err)
			assert.Equal(t, http.StatusNotImplemented, resp.StatusCode)
		})
	})

	t.Run("/v1/sensor/:id/data", func(t *testing.T) {
		t.Run("should call GET handler", func(t *testing.T) {
			mockSensorService := serviceMock.NewMockSensorService(t)
			app := sensor.RegisterRoutes(mockSensorService)

			// when
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/1/data", nil)

			// then
			resp, err := app.Test(req)
			defer resp.Body.Close()
			assert.NoError(t, err)
			assert.Equal(t, http.StatusNotImplemented, resp.StatusCode)
		})
	})
}

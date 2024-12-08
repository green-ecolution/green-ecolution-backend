package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	serviceMock "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	app := fiber.New()
  svc := serviceMock.NewMockServicesInterface(t)
  svc.EXPECT().AllServicesReady().Return(true)
  handler := HealthCheck(svc)
	app.Use(handler)

	t.Run("should return 200 OK for liveness probe", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should return 200 OK for readiness probe when services are ready", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ready", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should return 404 for undefined endpoint", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/undefined", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestHealthCheck_EdgeCases(t *testing.T) {
	app := fiber.New()

	handler := healthcheck.New(healthcheck.Config{
		LivenessProbe: func(_ *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/health",
		ReadinessProbe: func(_ *fiber.Ctx) bool {
			return false
		},
		ReadinessEndpoint: "/ready",
	})
	app.Use(handler)

	t.Run("should return 503 Service Unavailable when readiness probe fails", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ready", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	})

	t.Run("should handle high frequency of requests on /health", func(t *testing.T) {
		const requestCount = 100
		results := make(chan int, requestCount)

		for i := 0; i < requestCount; i++ {
			go func() {
				req := httptest.NewRequest(http.MethodGet, "/health", nil)
				resp, err := app.Test(req, -1)
				if err == nil {
					results <- resp.StatusCode
				} else {
					results <- http.StatusInternalServerError
				}
			}()
		}

		for i := 0; i < requestCount; i++ {
			status := <-results
			assert.Equal(t, http.StatusOK, status)
		}
	})

	t.Run("should reject POST request to /health", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/health", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should handle request with unusual headers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		req.Header.Set("X-Custom-Header", "CustomValue")
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should handle long URL paths on /health", func(t *testing.T) {
		longPath := "/health/" + strings.Repeat("pathsegment/", 50)
		req := httptest.NewRequest(http.MethodGet, longPath, nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should handle request with unsupported MIME type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		req.Header.Set("Content-Type", "application/xml")
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

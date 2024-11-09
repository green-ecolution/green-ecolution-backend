package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/middleware"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckMiddleware(t *testing.T) {
	t.Run("should return status 200 for health check", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.HealthCheck(nil))

		req := httptest.NewRequest("GET", "/health", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Health check should return status 200")
	})
	t.Run("should return status 200 for valid health check request", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.HealthCheck(nil))

		req := httptest.NewRequest("GET", "/health", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Health check should return status 200")
	})

	t.Run("should handle requests with query parameters gracefully", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.HealthCheck(nil))

		req := httptest.NewRequest("GET", "/health?param=value", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Health check should handle query parameters")
	})

	t.Run("should handle health check with different HTTP methods", func(t *testing.T) {
		methods := []string{"POST", "PUT", "DELETE", "PATCH"}
		for _, method := range methods {
			t.Run("method: "+method, func(t *testing.T) {
				app := fiber.New()
				app.Use(middleware.HealthCheck(nil))

				req := httptest.NewRequest(method, "/health", nil)
				resp, err := app.Test(req, -1)
				defer resp.Body.Close()

				assert.Nil(t, err)
				assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Health check should only allow GET method")
			})
		}
	})

	t.Run("should return proper headers in response", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.HealthCheck(nil))

		req := httptest.NewRequest("GET", "/health", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Health check should return status 200")
		assert.Equal(t, "text/plain; charset=utf-8", resp.Header.Get("Content-Type"), "Content-Type should be set to text/plain")
	})

	t.Run(
		"should return correct response body",
		func(t *testing.T) {
			app := fiber.New()
			app.Use(middleware.HealthCheck(nil))

			req := httptest.NewRequest("GET", "/health", nil)
			resp, err := app.Test(req, -1)
			defer resp.Body.Close()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode, "Health check should return status 200")

			body := make([]byte, resp.ContentLength)
			resp.Body.Read(body)
			assert.Equal(t, "OK", string(body), "Health check response body should be 'OK'")
		},
	)
}

func TestHTTPLoggerMiddleware(t *testing.T) {
	t.Run("should log and return a response", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.HTTPLogger())

		req := httptest.NewRequest("GET", "/log", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		assert.Nil(t, err)
		assert.NotNil(t, resp, "Response should not be nil for HTTPLogger")
	})

	t.Run("should log and return a response for valid GET request", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.HTTPLogger())

		app.Get("/log", func(c *fiber.Ctx) error {
			return c.SendString("Log endpoint")
		})

		req := httptest.NewRequest("GET", "/log", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Response should return status 200")
	})

	t.Run("should return 404 for non-existent route with logging", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.HTTPLogger())

		req := httptest.NewRequest("GET", "/non-existent", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Response should return status 404 for non-existent route")
	})

	t.Run("should handle logging for POST requests", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.HTTPLogger())

		app.Post("/log", func(c *fiber.Ctx) error {
			return c.SendString("Logged POST request")
		})

		req := httptest.NewRequest("POST", "/log", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "POST request should return status 200")
	})

	t.Run("should handle requests with headers and log them", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.HTTPLogger())

		app.Get("/log", func(c *fiber.Ctx) error {
			return c.SendString("Log endpoint with headers")
		})

		req := httptest.NewRequest("GET", "/log", nil)
		req.Header.Set("X-Custom-Header", "TestHeader")
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Request with headers should return status 200")

	})

	t.Run("should handle large payloads and log request details", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			BodyLimit: 20 * 1024 * 1024,
		})
		app.Use(middleware.HTTPLogger())

		app.Post("/log", func(c *fiber.Ctx) error {
			return c.SendString("Handled large payload")
		})

		largePayload := make([]byte, 10*1024*1024)
		req := httptest.NewRequest("POST", "/log", bytes.NewBuffer(largePayload))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		if resp != nil {
			defer resp.Body.Close()
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Large payload should be handled successfully")
	})

	t.Run("should handle 500 Internal Server Error response and log it", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.HTTPLogger())

		app.Get("/error", func(c *fiber.Ctx) error {
			return fiber.NewError(http.StatusInternalServerError, "Internal error")
		})

		req := httptest.NewRequest("GET", "/error", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "Internal server error should return status 500")
	})

	t.Run("should log requests with query parameters", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.HTTPLogger())

		app.Get("/log", func(c *fiber.Ctx) error {
			return c.SendString("Logged request with query parameters")
		})

		req := httptest.NewRequest("GET", "/log?param1=value1&param2=value2", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Request with query parameters should return status 200")
	})

}

func TestRequestIDMiddleware(t *testing.T) {
	t.Run("should set X-Request-ID header", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.RequestID())
		app.Get("/request-id", func(c *fiber.Ctx) error {
			return c.SendStatus(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "/request-id", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "RequestID middleware should allow the request to proceed")
		assert.NotEmpty(t, resp.Header.Get("X-Request-ID"), "X-Request-ID header should be set")
	})
}

func TestPublicRoutes(t *testing.T) {
	t.Run("should handle public route correctly", func(t *testing.T) {
		app := fiber.New()
		initPublicRoutes := func(app *fiber.App) {
			app.Get("/public", func(c *fiber.Ctx) error {
				return c.SendString("Public Route")
			})
		}

		app.Use(middleware.HealthCheck(nil))
		app.Use(middleware.RequestID())
		initPublicRoutes(app)

		req := httptest.NewRequest("GET", "/public", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Public route should return status 200")
	})
}

func TestInvalidRoutes(t *testing.T) {
	t.Run("should return 404 for invalid route", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.HealthCheck(nil))

		req := httptest.NewRequest("GET", "/invalid-route", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Invalid route should return status 404")
	})
}

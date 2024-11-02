package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestHTTPLogger(t *testing.T) {
	t.Run("should log GET request successfully", func(t *testing.T) {
		app := fiber.New()
		app.Use(HTTPLogger())

		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendStatus(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should log POST request with JSON payload", func(t *testing.T) {
		app := fiber.New()
		app.Use(HTTPLogger())

		app.Post("/test", func(c *fiber.Ctx) error {
			return c.SendStatus(http.StatusCreated)
		})

		payload := []byte(`{"name": "test"}`)
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("should log PUT request with large payload", func(t *testing.T) {
		app := fiber.New()
		app.Use(HTTPLogger())

		app.Put("/test", func(c *fiber.Ctx) error {
			return c.SendStatus(http.StatusOK)
		})

		largePayload := make([]byte, 1024*1024) // 1MB payload
		req := httptest.NewRequest(http.MethodPut, "/test", bytes.NewReader(largePayload))
		req.Header.Set("Content-Type", "application/octet-stream")
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should log DELETE request successfully", func(t *testing.T) {
		app := fiber.New()
		app.Use(HTTPLogger())

		app.Delete("/test", func(c *fiber.Ctx) error {
			return c.SendStatus(http.StatusNoContent)
		})

		req := httptest.NewRequest(http.MethodDelete, "/test", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("should log request with query parameters", func(t *testing.T) {
		app := fiber.New()
		app.Use(HTTPLogger())

		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString(c.Query("param"))
		})

		req := httptest.NewRequest(http.MethodGet, "/test?param=value", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should handle request with missing route gracefully", func(t *testing.T) {
		app := fiber.New()
		app.Use(HTTPLogger())

		req := httptest.NewRequest(http.MethodGet, "/missing", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should handle HEAD request", func(t *testing.T) {
		app := fiber.New()
		app.Use(HTTPLogger())

		app.Head("/test", func(c *fiber.Ctx) error {
			return c.SendStatus(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodHead, "/test", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

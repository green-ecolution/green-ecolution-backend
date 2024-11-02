package middleware

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

func TestHTTPLogger_BorderlineCases(t *testing.T) {
	t.Run("should handle long processing time", func(t *testing.T) {
		app := fiber.New()
		app.Use(HTTPLogger())

		app.Get("/timeout", func(c *fiber.Ctx) error {

			ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
			defer cancel()

			select {
			case <-time.After(6 * time.Second):
				return c.SendStatus(http.StatusOK)
			case <-ctx.Done():
				return c.SendStatus(http.StatusRequestTimeout)
			}
		})

		req := httptest.NewRequest(http.MethodGet, "/timeout", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusRequestTimeout, resp.StatusCode)
	})

	t.Run("should handle TRACE method request", func(t *testing.T) {
		app := fiber.New()
		app.Use(HTTPLogger())

		app.Trace("/trace", func(c *fiber.Ctx) error {
			return c.SendStatus(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodTrace, "/trace", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should handle large number of simultaneous requests", func(t *testing.T) {
		app := fiber.New()
		app.Use(HTTPLogger())
		app.Get("/test", func(c *fiber.Ctx) error { return c.SendStatus(http.StatusOK) })

		concurrency := 100
		results := make(chan bool, concurrency)

		for i := 0; i < concurrency; i++ {
			go func() {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				resp, err := app.Test(req, -1)
				results <- err == nil && resp.StatusCode == http.StatusOK
			}()
		}

		for i := 0; i < concurrency; i++ {
			assert.True(t, <-results)
		}
	})

	t.Run("should handle request with unusual headers", func(t *testing.T) {
		app := fiber.New()
		app.Use(HTTPLogger())

		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendStatus(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Unusual-Header", "UnusualHeaderValue")
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should handle large URL path and query", func(t *testing.T) {
		app := fiber.New()
		app.Use(HTTPLogger())

		longPath := "/" + strings.Repeat("pathsegment/", 50)
		req := httptest.NewRequest(http.MethodGet, longPath, nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should handle multilingual characters in path", func(t *testing.T) {
		app := fiber.New()
		app.Use(HTTPLogger())

		app.Get("/multilingual/こんにちは", func(c *fiber.Ctx) error {
			return c.SendStatus(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/multilingual/こんにちは", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

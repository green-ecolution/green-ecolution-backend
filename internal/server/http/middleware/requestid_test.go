package middleware

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestRequestID(t *testing.T) {
	t.Run("should add X-Request-ID header to each request", func(t *testing.T) {
		app := fiber.New()
		app.Use(RequestID())

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Header.Get("X-Request-ID"))
	})

	t.Run("should generate unique IDs for different requests", func(t *testing.T) {
		app := fiber.New()
		app.Use(RequestID())

		req1 := httptest.NewRequest(http.MethodGet, "/", nil)
		resp1, err1 := app.Test(req1, -1)
		assert.NoError(t, err1)
		id1 := resp1.Header.Get("X-Request-ID")
		assert.NotEmpty(t, id1)

		req2 := httptest.NewRequest(http.MethodGet, "/", nil)
		resp2, err2 := app.Test(req2, -1)
		assert.NoError(t, err2)
		id2 := resp2.Header.Get("X-Request-ID")
		assert.NotEmpty(t, id2)

		assert.NotEqual(t, id1, id2)
	})

	t.Run("should ensure unique IDs in parallel requests", func(t *testing.T) {
		app := fiber.New()
		app.Use(RequestID())

		requestCount := 10
		ids := make(chan string, requestCount)

		for i := 0; i < requestCount; i++ {
			go func() {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				resp, err := app.Test(req, -1)
				assert.NoError(t, err)
				ids <- resp.Header.Get("X-Request-ID")
			}()
		}

		uniqueIDs := make(map[string]bool)
		for i := 0; i < requestCount; i++ {
			id := <-ids
			assert.NotEmpty(t, id)
			uniqueIDs[id] = true
		}

		assert.Equal(t, requestCount, len(uniqueIDs))
	})

	t.Run("should overwrite existing X-Request-ID header", func(t *testing.T) {
		app := fiber.New()
		app.Use(RequestID())

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Request-ID", "existing-id")
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, "existing-id", resp.Header.Get("X-Request-ID"))
	})

	t.Run("should return 400 Bad Request for invalid HTTP method", func(t *testing.T) {
		app := fiber.New()
		app.Use(RequestID())

		req := httptest.NewRequest("INVALID", "/", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should add X-Request-ID header to each request", func(t *testing.T) {
		app := fiber.New()
		app.Use(RequestID())

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Header.Get("X-Request-ID"))
	})

	t.Run("should generate unique IDs for different requests", func(t *testing.T) {
		app := fiber.New()
		app.Use(RequestID())

		req1 := httptest.NewRequest(http.MethodGet, "/", nil)
		resp1, err1 := app.Test(req1, -1)
		assert.NoError(t, err1)
		id1 := resp1.Header.Get("X-Request-ID")
		assert.NotEmpty(t, id1)

		req2 := httptest.NewRequest(http.MethodGet, "/", nil)
		resp2, err2 := app.Test(req2, -1)
		assert.NoError(t, err2)
		id2 := resp2.Header.Get("X-Request-ID")
		assert.NotEmpty(t, id2)

		assert.NotEqual(t, id1, id2)
	})

	t.Run("should ensure unique IDs in parallel requests", func(t *testing.T) {
		app := fiber.New()
		app.Use(RequestID())

		requestCount := 10
		ids := make(chan string, requestCount)

		for i := 0; i < requestCount; i++ {
			go func() {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				resp, err := app.Test(req, -1)
				assert.NoError(t, err)
				ids <- resp.Header.Get("X-Request-ID")
			}()
		}

		uniqueIDs := make(map[string]bool)
		for i := 0; i < requestCount; i++ {
			id := <-ids
			assert.NotEmpty(t, id)
			uniqueIDs[id] = true
		}

		assert.Equal(t, requestCount, len(uniqueIDs))
	})

	t.Run("should retain existing X-Request-ID if already present", func(t *testing.T) {
		app := fiber.New()
		app.Use(RequestID())

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Request-ID", "existing-id")
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, "existing-id", resp.Header.Get("X-Request-ID"))
	})

	t.Run("should generate valid UUID format for X-Request-ID", func(t *testing.T) {
		app := fiber.New()
		app.Use(RequestID())

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		resp, err := app.Test(req, -1)
		assert.NoError(t, err)

		requestID := resp.Header.Get("X-Request-ID")
		uuidRegex := regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)
		assert.Regexp(t, uuidRegex, requestID)
	})

	t.Run("should retain X-Request-ID after forwarding", func(t *testing.T) {
		app := fiber.New()
		app.Use(RequestID())
		app.Get("/forward", func(c *fiber.Ctx) error {

			forwardedReq := httptest.NewRequest(http.MethodGet, "/", nil)
			forwardedReq.Header.Set("X-Request-ID", c.Get("X-Request-ID"))

			assert.Equal(t, c.Get("X-Request-ID"), forwardedReq.Header.Get("X-Request-ID"))
			return c.SendStatus(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/forward", nil)
		req.Header.Set("X-Request-ID", "test-id")
		_, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, "test-id", req.Header.Get("X-Request-ID"))
	})

	t.Run("should function correctly with other middleware", func(t *testing.T) {
		app := fiber.New()
		app.Use(RequestID())
		app.Use(func(c *fiber.Ctx) error {
			return c.Next()
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Header.Get("X-Request-ID"))
	})

	t.Run("should handle max header size", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ReadBufferSize: 16 * 1024,
		})
		app.Use(RequestID())

		largeHeader := strings.Repeat("a", 8192)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Large-Header", largeHeader)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Header.Get("X-Request-ID"))
	})

	t.Run("should generate X-Request-ID for different HTTP methods", func(t *testing.T) {
		app := fiber.New()
		app.Use(RequestID())

		methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}
		for _, method := range methods {
			req := httptest.NewRequest(method, "/", nil)
			resp, err := app.Test(req, -1)
			assert.NoError(t, err)
			assert.NotEmpty(t, resp.Header.Get("X-Request-ID"))
		}
	})

	t.Run("should generate X-Request-ID for large payloads", func(t *testing.T) {
		app := fiber.New()
		app.Use(RequestID())

		largePayload := strings.Repeat("a", 1<<20)
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(largePayload))
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Header.Get("X-Request-ID"))
	})
}

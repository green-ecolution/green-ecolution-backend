package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// Test 1: Check if X-Request-ID header is present
func TestRequestID(t *testing.T) {
	app := fiber.New()
	app.Use(RequestID())

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Header.Get("X-Request-ID"))
}

// Test 2: Check if different requests receive unique IDs
func TestRequestID_UniqueIDs(t *testing.T) {
	app := fiber.New()
	app.Use(RequestID())

	// First request
	req1 := httptest.NewRequest("GET", "/", nil)
	resp1, err1 := app.Test(req1)
	assert.NoError(t, err1)
	id1 := resp1.Header.Get("X-Request-ID")
	assert.NotEmpty(t, id1)

	// Second request
	req2 := httptest.NewRequest("GET", "/", nil)
	resp2, err2 := app.Test(req2)
	assert.NoError(t, err2)
	id2 := resp2.Header.Get("X-Request-ID")
	assert.NotEmpty(t, id2)

	assert.NotEqual(t, id1, id2)
}

// Test 3: Parallel requests
func TestRequestID_ParallelRequests(t *testing.T) {
	app := fiber.New()
	app.Use(RequestID())

	requestCount := 10
	ids := make(chan string, requestCount)

	for i := 0; i < requestCount; i++ {
		go func() {
			req := httptest.NewRequest("GET", "/", nil)
			resp, err := app.Test(req)
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
}

// Test 4: Invalid request method
func TestRequestID_InvalidMethod(t *testing.T) {
	app := fiber.New()
	app.Use(RequestID())

	req := httptest.NewRequest("INVALID", "/", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

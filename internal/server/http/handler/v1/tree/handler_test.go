package tree

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	httpEntities "github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	serviceMock "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestGetAllTrees tests the GetAllTrees handler.
func TestGetAllTrees(t *testing.T) {
	t.Run("should return all trees successfully", func(t *testing.T) {
		app := fiber.New()
		mockService := serviceMock.NewMockTreeService(t)
		app.Get("/v1/tree", GetAllTrees(mockService))

		mockService.EXPECT().GetAll(
			mock.Anything,
		).Return(getMockTrees(), nil)

		// Create a request to the endpoint
		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/tree", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// Verify the response
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
	t.Run("should return an empty list when no trees are available", func(t *testing.T) {
		app := fiber.New()
		mockService := serviceMock.NewMockTreeService(t)
		app.Get("/v1/tree", GetAllTrees(mockService))

		mockService.EXPECT().GetAll(
			mock.Anything,
		).Return([]*entities.Tree{}, nil)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/tree", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
	t.Run("should return 500 when internal server error occurs", func(t *testing.T) {
		app := fiber.New()
		mockService := serviceMock.NewMockTreeService(t)
		app.Get("/v1/tree", GetAllTrees(mockService))

		mockService.EXPECT().GetAll(
			mock.Anything,
		).Return(nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error"))

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/tree", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockService.AssertExpectations(t)
	})
}

// TestCreateTree tests the CreateTree handler.
func TestCreateTree(t *testing.T) {
	t.Run("should create tree cluster successfully", func(t *testing.T) {
		app := fiber.New()
		mockService := serviceMock.NewMockTreeService(t)
		app.Post("/v1/tree", CreateTree(mockService))

		testTree := getMockTrees()[0]

		mockService.EXPECT().Create(
			mock.Anything,
			mock.AnythingOfType("*entities.TreeCreate"),
		).Return(testTree, nil)

		// when
		reqBody := getMockTreeCreateRequest()
		reqBodyBytes, _ := json.Marshal(reqBody)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/v1/tree", bytes.NewBuffer(reqBodyBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response httpEntities.TreeResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, testTree.Latitude, response.Latitude)
		assert.Equal(t, testTree.Longitude, response.Longitude)

		mockService.AssertExpectations(t)
	})

	t.Run("should return 400 Bad Request for invalid request body", func(t *testing.T) {
		app := fiber.New()
		mockService := serviceMock.NewMockTreeService(t)
		app.Post("/v1/tree", CreateTree(mockService))

		invalidRequestBody := []byte(`{"invalid_field": "value"}`)

		// when
		body, _ := json.Marshal(invalidRequestBody)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/v1/tree", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
	t.Run("should return 500 when internal server error occurs", func(t *testing.T) {
		app := fiber.New()
		mockService := serviceMock.NewMockTreeService(t)
		app.Post("/v1/tree", CreateTree(mockService))

		mockService.EXPECT().Create(
			mock.Anything,
			mock.AnythingOfType("*entities.TreeCreate"),
		).Return(nil, fiber.NewError(fiber.StatusInternalServerError, "service error"))

		// when
		body, _ := json.Marshal(getMockTreeCreateRequest())
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/v1/tree", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}

func getMockTreeCreateRequest() *httpEntities.TreeCreateRequest {
	return &httpEntities.TreeCreateRequest{
		TreeClusterID: nil,
		Readonly:      false,
		PlantingYear:  2023,
		Species:       "Oak",
		Number:        "T001",
		Latitude:      54.801539,
		Longitude:     9.446741,
		SensorID:      nil,
		Description:   "A newly planted oak tree",
	}
}
func getMockTrees() []*entities.Tree {
	now := time.Now()

	return []*entities.Tree{
		{
			ID:           1,
			CreatedAt:    now,
			UpdatedAt:    now,
			Species:      "Oak",
			Number:       "T001",
			Latitude:     9.446741,
			Longitude:    54.801539,
			Description:  "A mature oak tree",
			PlantingYear: 2023,
			Readonly:     true,
		},
		{
			ID:           2,
			CreatedAt:    now,
			UpdatedAt:    now,
			Species:      "Pine",
			Number:       "T002",
			Latitude:     9.446700,
			Longitude:    54.801510,
			Description:  "A young pine tree",
			PlantingYear: 2023,
			Readonly:     true,
		},
	}
}

package tree

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	serviceMock "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
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

func getMockTrees() []*entities.Tree {
	return []*entities.Tree{
		{
			ID:           1,
			Species:      "Oak",
			PlantingYear: 2023,
			Number:       "T001",
			Latitude:     54.801539,
			Longitude:    9.446741,
			Description:  "oak tree",
		},
		{
			ID:           2,
			Species:      "Pine",
			PlantingYear: 2022,
			Number:       "T002",
			Latitude:     54.801500,
			Longitude:    9.446700,
			Description:  "pine tree",
		},
	}
}

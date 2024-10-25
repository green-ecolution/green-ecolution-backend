package treecluster

import (
	"context"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	serverEntities "github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	serviceMock "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllTreeClusters(t *testing.T) {
	t.Run("should return all tree clusters successfully", func(t *testing.T) {
		app := fiber.New()
		mockClusterService := serviceMock.NewMockTreeClusterService(t)
		handler := GetAllTreeClusters(mockClusterService)
		app.Get("/v1/cluster", handler)

		expectedData := []*entities.TreeCluster{
			{ID: 1, Name: "Cluster A"},
			{ID: 2, Name: "Cluster B"},
		}

		mockClusterService.EXPECT().GetAll(mock.Anything).Return(expectedData, nil)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), "GET", "/v1/cluster", nil)
		resp, err := app.Test(req, -1)

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response serverEntities.TreeClusterListResponse
		err = utils.ParseJSONResponse(resp, &response)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(response.Data))
		assert.Equal(t, "Cluster A", response.Data[0].Name)

		mockClusterService.AssertExpectations(t)
	})

	t.Run("should return an empty list when no tree clusters are available", func(t *testing.T) {
		app := fiber.New()
		mockClusterService := serviceMock.NewMockTreeClusterService(t)
		handler := GetAllTreeClusters(mockClusterService)
		app.Get("/v1/cluster", handler)

		expectedData := []*entities.TreeCluster{}

		mockClusterService.EXPECT().GetAll(mock.Anything).Return(expectedData, nil)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), "GET", "/v1/cluster", nil)
		resp, err := app.Test(req, -1)

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response serverEntities.TreeClusterListResponse
		err = utils.ParseJSONResponse(resp, &response)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(response.Data))

		mockClusterService.AssertExpectations(t)
	})

	t.Run("should return 500 Internal Server Error when service fails", func(t *testing.T) {
		app := fiber.New()
		mockClusterService := serviceMock.NewMockTreeClusterService(t)
		handler := GetAllTreeClusters(mockClusterService)
		app.Get("/v1/cluster", handler)

		mockClusterService.EXPECT().GetAll(mock.Anything).Return(nil, fiber.NewError(fiber.StatusInternalServerError, "service error"))

		// when
		req, _ := http.NewRequestWithContext(context.Background(), "GET", "/v1/cluster", nil)
		resp, err := app.Test(req, -1)

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockClusterService.AssertExpectations(t)
	})
}

func TestGetTreeClusterByID(t *testing.T) {
	t.Run("should return tree cluster successfully", func(t *testing.T) {
		app := fiber.New()
		mockClusterService := serviceMock.NewMockTreeClusterService(t)
		handler := GetTreeClusterByID(mockClusterService)
		app.Get("/v1/cluster/:treecluster_id", handler)

		expectedData := &entities.TreeCluster{
			ID:   1,
			Name: "Cluster A",
		}

		mockClusterService.EXPECT().GetByID(mock.Anything, int32(1)).Return(expectedData, nil)

		req, _ := http.NewRequestWithContext(context.Background(), "GET", "/v1/cluster/1", nil)
		resp, err := app.Test(req, -1)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response serverEntities.TreeClusterResponse
		err = utils.ParseJSONResponse(resp, &response)
		assert.NoError(t, err)
		assert.Equal(t, "Cluster A", response.Name)

		mockClusterService.AssertExpectations(t)
	})

	t.Run("should return 400 Bad Request for invalid cluster ID", func(t *testing.T) {
		app := fiber.New()
		mockClusterService := serviceMock.NewMockTreeClusterService(t)
		handler := GetTreeClusterByID(mockClusterService)
		app.Get("/v1/cluster/:treecluster_id", handler)

		req, _ := http.NewRequestWithContext(context.Background(), "GET", "/v1/cluster/invalid", nil)
		resp, err := app.Test(req, -1)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 404 Not Found if cluster does not exist", func(t *testing.T) {
		app := fiber.New()
		mockClusterService := serviceMock.NewMockTreeClusterService(t)
		handler := GetTreeClusterByID(mockClusterService)
		app.Get("/v1/cluster/:treecluster_id", handler)

		mockClusterService.EXPECT().GetByID(mock.Anything, int32(999)).Return(nil, storage.ErrTreeClusterNotFound)

		req, _ := http.NewRequestWithContext(context.Background(), "GET", "/v1/cluster/999", nil)
		resp, err := app.Test(req, -1)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		mockClusterService.AssertExpectations(t)
	})

	t.Run("should return 500 Internal Server Error for service failure", func(t *testing.T) {
		app := fiber.New()
		mockClusterService := serviceMock.NewMockTreeClusterService(t)
		handler := GetTreeClusterByID(mockClusterService)
		app.Get("/v1/cluster/:treecluster_id", handler)

		mockClusterService.EXPECT().GetByID(mock.Anything, int32(1)).Return(nil, fiber.NewError(fiber.StatusInternalServerError, "service error"))

		req, _ := http.NewRequestWithContext(context.Background(), "GET", "/v1/cluster/1", nil)
		resp, err := app.Test(req, -1)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockClusterService.AssertExpectations(t)
	})
}

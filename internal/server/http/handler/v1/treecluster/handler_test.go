package treecluster_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	serverEntities "github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/treecluster"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	serviceMock "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllTreeCluster(t *testing.T) {
	t.Run("should return all tree clusters successfully", func(t *testing.T) {
		app := fiber.New()
		mockClusterService := serviceMock.NewMockTreeClusterService(t)
		handler := treecluster.GetAllTreeClusters(mockClusterService)
		app.Get("/v1/cluster", handler)

		mockClusterService.EXPECT().GetAll(
			mock.Anything,
		).Return(TestClusterList, nil)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/cluster", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response serverEntities.TreeClusterListResponse
		err = utils.ParseJSONResponse(resp, &response)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(response.Data))
		assert.Equal(t, TestClusterList[0].Name, response.Data[0].Name)

		mockClusterService.AssertExpectations(t)
	})

	t.Run("should return an empty list when no tree clusters are available", func(t *testing.T) {
		app := fiber.New()
		mockClusterService := serviceMock.NewMockTreeClusterService(t)
		handler := treecluster.GetAllTreeClusters(mockClusterService)
		app.Get("/v1/cluster", handler)

		mockClusterService.EXPECT().GetAll(
			mock.Anything,
		).Return([]*entities.TreeCluster{}, nil)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/cluster", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

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
		handler := treecluster.GetAllTreeClusters(mockClusterService)
		app.Get("/v1/cluster", handler)

		mockClusterService.EXPECT().GetAll(
			mock.Anything,
		).Return(nil, fiber.NewError(fiber.StatusInternalServerError, "service error"))

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/cluster", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

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
		handler := treecluster.GetTreeClusterByID(mockClusterService)
		app.Get("/v1/cluster/:treecluster_id", handler)

		mockClusterService.EXPECT().GetByID(
			mock.Anything,
			int32(1),
		).Return(TestCluster, nil)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/cluster/1", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response serverEntities.TreeClusterResponse
		err = utils.ParseJSONResponse(resp, &response)
		assert.NoError(t, err)
		assert.Equal(t, TestCluster.Name, response.Name)

		mockClusterService.AssertExpectations(t)
	})

	t.Run("should return 400 Bad Request for invalid cluster ID", func(t *testing.T) {
		app := fiber.New()
		mockClusterService := serviceMock.NewMockTreeClusterService(t)
		handler := treecluster.GetTreeClusterByID(mockClusterService)
		app.Get("/v1/cluster/:treecluster_id", handler)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/cluster/invalid", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 404 Not Found if cluster does not exist", func(t *testing.T) {
		app := fiber.New()
		mockClusterService := serviceMock.NewMockTreeClusterService(t)
		handler := treecluster.GetTreeClusterByID(mockClusterService)
		app.Get("/v1/cluster/:treecluster_id", handler)

		mockClusterService.EXPECT().GetByID(
			mock.Anything,
			int32(999),
		).Return(nil, service.NewError(service.NotFound, "not found"))

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/cluster/999", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		mockClusterService.AssertExpectations(t)
	})

	t.Run("should return 500 Internal Server Error for service failure", func(t *testing.T) {
		app := fiber.New()
		mockClusterService := serviceMock.NewMockTreeClusterService(t)
		handler := treecluster.GetTreeClusterByID(mockClusterService)
		app.Get("/v1/cluster/:treecluster_id", handler)

		mockClusterService.EXPECT().GetByID(
			mock.Anything,
			int32(1),
		).Return(nil, fiber.NewError(fiber.StatusInternalServerError, "service error"))

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/cluster/1", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockClusterService.AssertExpectations(t)
	})
}

func TestCreateTreeCluster(t *testing.T) {
	t.Run("should create tree cluster successfully", func(t *testing.T) {
		app := fiber.New()
		mockClusterService := serviceMock.NewMockTreeClusterService(t)
		handler := treecluster.CreateTreeCluster(mockClusterService)
		app.Post("/v1/cluster", handler)

		mockClusterService.EXPECT().Create(
			mock.Anything,
			mock.AnythingOfType("*entities.TreeClusterCreate"),
		).Return(TestCluster, nil)

		// when
		body, _ := json.Marshal(TestClusterRequest)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/v1/cluster", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response serverEntities.TreeClusterResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, TestCluster.Name, response.Name)
		assert.Equal(t, TestCluster.Region.Name, response.Region.Name)

		mockClusterService.AssertExpectations(t)
	})

	t.Run("should return 400 Bad Request for invalid request body", func(t *testing.T) {
		app := fiber.New()
		mockClusterService := serviceMock.NewMockTreeClusterService(t)
		handler := treecluster.CreateTreeCluster(mockClusterService)
		app.Post("/v1/cluster", handler)

		invalidRequestBody := []byte(`{"invalid_field": "value"}`)

		// when
		body, _ := json.Marshal(invalidRequestBody)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/v1/cluster", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 500 Internal Server Error for service failure", func(t *testing.T) {
		app := fiber.New()
		mockClusterService := serviceMock.NewMockTreeClusterService(t)
		handler := treecluster.CreateTreeCluster(mockClusterService)
		app.Post("/v1/cluster", handler)

		mockClusterService.EXPECT().Create(
			mock.Anything,
			mock.AnythingOfType("*entities.TreeClusterCreate"),
		).Return(nil, fiber.NewError(fiber.StatusInternalServerError, "service error"))

		// when
		body, _ := json.Marshal(TestClusterRequest)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/v1/cluster", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockClusterService.AssertExpectations(t)
	})
}

func TestUpdateTreeCluster(t *testing.T) {
	t.Run("should update tree cluster successfully", func(t *testing.T) {
		app := fiber.New()
		mockClusterService := serviceMock.NewMockTreeClusterService(t)
		handler := treecluster.UpdateTreeCluster(mockClusterService)
		app.Put("/v1/cluster/:treecluster_id", handler)

		mockClusterService.EXPECT().Update(
			mock.Anything,
			int32(1),
			mock.Anything,
		).Return(TestCluster, nil)

		// when
		body, _ := json.Marshal(TestClusterRequest)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, "/v1/cluster/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response serverEntities.TreeClusterResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, TestCluster.Name, response.Name)

		mockClusterService.AssertExpectations(t)
	})

	t.Run("should return 400 Bad Request for invalid cluster ID", func(t *testing.T) {
		app := fiber.New()
		mockClusterService := serviceMock.NewMockTreeClusterService(t)
		handler := treecluster.UpdateTreeCluster(mockClusterService)
		app.Put("/v1/cluster/:treecluster_id", handler)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, "/v1/cluster/invalid", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 Bad Request for invalid request body", func(t *testing.T) {
		app := fiber.New()
		mockClusterService := serviceMock.NewMockTreeClusterService(t)
		handler := treecluster.UpdateTreeCluster(mockClusterService)
		app.Put("/v1/cluster/:treecluster_id", handler)

		invalidRequestBody := []byte(`{"invalid_field": "value"}`)

		// when
		body, _ := json.Marshal(invalidRequestBody)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, "/v1/cluster/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 404 Not Found if cluster does not exist", func(t *testing.T) {
		app := fiber.New()
		mockClusterService := serviceMock.NewMockTreeClusterService(t)
		handler := treecluster.UpdateTreeCluster(mockClusterService)
		app.Put("/v1/cluster/:treecluster_id", handler)

		mockClusterService.EXPECT().Update(
			mock.Anything,
			int32(1),
			mock.Anything,
		).Return(nil, service.NewError(service.NotFound, "not found"))

		// when
		body, _ := json.Marshal(TestClusterRequest)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, "/v1/cluster/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		mockClusterService.AssertExpectations(t)
	})

	t.Run("should return 500 Internal Server Error for service failure", func(t *testing.T) {
		app := fiber.New()
		mockClusterService := serviceMock.NewMockTreeClusterService(t)
		handler := treecluster.UpdateTreeCluster(mockClusterService)
		app.Put("/v1/cluster/:treecluster_id", handler)

		mockClusterService.EXPECT().Update(mock.Anything, int32(1), mock.Anything).Return(nil, fiber.NewError(fiber.StatusInternalServerError, "service error"))

		// when
		body, _ := json.Marshal(TestClusterRequest)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, "/v1/cluster/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockClusterService.AssertExpectations(t)
	})
}

func TestDeleteTreeCluster(t *testing.T) {
	t.Run("should delete tree cluster successfully", func(t *testing.T) {
		app := fiber.New()
		mockClusterService := serviceMock.NewMockTreeClusterService(t)
		handler := treecluster.DeleteTreeCluster(mockClusterService)
		app.Delete("/v1/cluster/:treecluster_id", handler)

		clusterID := 1
		mockClusterService.EXPECT().Delete(mock.Anything, int32(clusterID)).Return(nil)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, "/v1/cluster/"+strconv.Itoa(clusterID), nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		mockClusterService.AssertExpectations(t)
	})

	t.Run("should return 400 for invalid ID format", func(t *testing.T) {
		app := fiber.New()
		mockClusterService := serviceMock.NewMockTreeClusterService(t)
		handler := treecluster.DeleteTreeCluster(mockClusterService)
		app.Delete("/v1/cluster/:treecluster_id", handler)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, "/v1/cluster/invalid_id", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 404 for non-existing tree cluster", func(t *testing.T) {
		app := fiber.New()
		mockClusterService := serviceMock.NewMockTreeClusterService(t)
		handler := treecluster.DeleteTreeCluster(mockClusterService)
		app.Delete("/v1/cluster/:treecluster_id", handler)

		clusterID := 999
		mockClusterService.EXPECT().Delete(
			mock.Anything,
			int32(clusterID),
		).Return(service.NewError(service.NotFound, "not found"))

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, "/v1/cluster/"+strconv.Itoa(clusterID), nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockClusterService.AssertExpectations(t)
	})
}

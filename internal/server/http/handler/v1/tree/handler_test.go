package tree_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	httpEntities "github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/tree"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	serviceMock "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllTrees(t *testing.T) {
	t.Run("should return all trees successfully", func(t *testing.T) {
		app := fiber.New()
		mockTreeService := serviceMock.NewMockTreeService(t)
		app.Get("/v1/tree", tree.GetAllTrees(mockTreeService))
		mockTreeService.EXPECT().GetAll(
			mock.Anything,
			"",
		).Return(TestTrees, nil)

		// when
		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/tree", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockTreeService.AssertExpectations(t)
	})

	t.Run("should return an empty list when no trees are available", func(t *testing.T) {
		app := fiber.New()
		mockTreeService := serviceMock.NewMockTreeService(t)
		app.Get("/v1/tree", tree.GetAllTrees(mockTreeService))
		mockTreeService.EXPECT().GetAll(
			mock.Anything,
			"",
		).Return([]*entities.Tree{}, nil)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/tree", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockTreeService.AssertExpectations(t)
	})

	t.Run("should return 500 when internal server error occurs", func(t *testing.T) {
		app := fiber.New()
		mockTreeService := serviceMock.NewMockTreeService(t)
		app.Get("/v1/tree", tree.GetAllTrees(mockTreeService))

		mockTreeService.EXPECT().GetAll(
			mock.Anything,
			"",
		).Return(nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error"))

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/tree", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockTreeService.AssertExpectations(t)
	})
}

func TestGetTreeBySensorID(t *testing.T) {
	t.Run("should return tree successfully", func(t *testing.T) {
		app := fiber.New()
		mockTreeService := serviceMock.NewMockTreeService(t)
		app.Get("v1/tree/sensor/:sensor_id", tree.GetTreeBySensorID(mockTreeService))

		sensorID := "sensor-1"
		mockTreeService.EXPECT().GetBySensorID(
			mock.Anything,
			sensorID,
		).Return(TestTrees[0], nil)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), "GET", "/v1/tree/sensor/"+sensorID, nil)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, -1)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockTreeService.AssertExpectations(t)
	})

	t.Run("should return 404 when tree not found", func(t *testing.T) {
		app := fiber.New()
		mockTreeService := serviceMock.NewMockTreeService(t)
		app.Get("v1/tree/sensor/:sensor_id", tree.GetTreeBySensorID(mockTreeService))

		sensorID := "sensor-999"
		mockTreeService.EXPECT().GetBySensorID(
			mock.Anything,
			sensorID,
		).Return(nil, service.NewError(service.NotFound, "not found"))

		req, _ := http.NewRequestWithContext(context.Background(), "GET", "/v1/tree/sensor/"+sensorID, nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockTreeService.AssertExpectations(t)
	})

	t.Run("should return 500 when internal server error occurs", func(t *testing.T) {
		app := fiber.New()
		mockTreeService := serviceMock.NewMockTreeService(t)
		app.Get("/v1/tree/sensor/:sensor_id", tree.GetTreeBySensorID(mockTreeService))

		sensorID := "sensor-1"
		mockTreeService.EXPECT().GetBySensorID(
			mock.Anything,
			sensorID,
		).Return(nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error"))

		req, _ := http.NewRequestWithContext(context.Background(), "GET", "/v1/tree/sensor/"+sensorID, nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockTreeService.AssertExpectations(t)
	})
}

func TestCreateTree(t *testing.T) {
	t.Run("should create tree cluster successfully", func(t *testing.T) {
		app := fiber.New()
		mockTreeService := serviceMock.NewMockTreeService(t)
		app.Post("/v1/tree", tree.CreateTree(mockTreeService))

		testTree := TestTrees[0]
		mockTreeService.EXPECT().Create(
			mock.Anything,
			mock.AnythingOfType("*entities.TreeCreate"),
		).Return(testTree, nil)

		// when
		reqBody := TestTreeCreateRequest
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

		mockTreeService.AssertExpectations(t)
	})

	t.Run("should return 400 Bad Request for invalid request body", func(t *testing.T) {
		app := fiber.New()
		mockTreeService := serviceMock.NewMockTreeService(t)
		app.Post("/v1/tree", tree.CreateTree(mockTreeService))
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
		mockTreeService := serviceMock.NewMockTreeService(t)
		app.Post("/v1/tree", tree.CreateTree(mockTreeService))
		mockTreeService.EXPECT().Create(
			mock.Anything,
			mock.AnythingOfType("*entities.TreeCreate"),
		).Return(nil, fiber.NewError(fiber.StatusInternalServerError, "service error"))

		// when
		body, _ := json.Marshal(TestTreeCreateRequest)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/v1/tree", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockTreeService.AssertExpectations(t)
	})
}

func TestUpdateTree(t *testing.T) {
	t.Run("should update tree successfully", func(t *testing.T) {
		app := fiber.New()
		mockTreeService := serviceMock.NewMockTreeService(t)
		app.Put("/v1/tree/:id", tree.UpdateTree(mockTreeService))
		testTree := TestTrees[0]
		treeID := int32(1)
		mockTreeService.EXPECT().Update(
			mock.Anything,
			treeID,
			mock.Anything,
		).Return(testTree, nil)

		// when
		body, _ := json.Marshal(TestTreeUpdateRequest)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, "/v1/tree/"+strconv.Itoa(int(treeID)), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response httpEntities.TreeResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, testTree.Latitude, response.Latitude)
		assert.Equal(t, testTree.Longitude, response.Longitude)

		mockTreeService.AssertExpectations(t)
	})

	t.Run("should return 400 Bad Request for invalid tree ID", func(t *testing.T) {
		app := fiber.New()
		mockTreeService := serviceMock.NewMockTreeService(t)
		handler := tree.UpdateTree(mockTreeService)
		app.Put("/v1/tree/:id", handler)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, "/v1/tree/invalid", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 Bad Request for invalid request body", func(t *testing.T) {
		app := fiber.New()
		mockTreeService := serviceMock.NewMockTreeService(t)
		app.Put("/v1/tree/:id", tree.UpdateTree(mockTreeService))

		invalidRequestBody := []byte(`{"invalid_field": "value"}`)
		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, "/v1/tree/1", bytes.NewBuffer(invalidRequestBody))
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 404 Not Found when tree does not exist", func(t *testing.T) {
		app := fiber.New()
		mockTreeService := serviceMock.NewMockTreeService(t)
		app.Put("/v1/tree/:id", tree.UpdateTree(mockTreeService))

		treeID := int32(999)
		mockTreeService.EXPECT().Update(
			mock.Anything,
			treeID,
			mock.AnythingOfType("*entities.TreeUpdate"),
		).Return(nil, service.NewError(service.NotFound, "not found"))

		// when
		reqBody := TestTreeUpdateRequest
		reqBodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/v1/tree/"+strconv.Itoa(int(treeID)), bytes.NewBuffer(reqBodyBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockTreeService.AssertExpectations(t)
	})

	t.Run("should return 500 Internal Server Error on service error", func(t *testing.T) {
		app := fiber.New()
		mockTreeService := serviceMock.NewMockTreeService(t)
		app.Put("/v1/tree/:id", tree.UpdateTree(mockTreeService))

		treeID := int32(1)
		mockTreeService.EXPECT().Update(
			mock.Anything,
			treeID,
			mock.AnythingOfType("*entities.TreeUpdate"),
		).Return(nil, fiber.NewError(fiber.StatusInternalServerError, "server error"))

		// when
		body, _ := json.Marshal(TestTreeUpdateRequest)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, "/v1/tree/"+strconv.Itoa(int(treeID)), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockTreeService.AssertExpectations(t)
	})
}

func TestDeleteTree(t *testing.T) {
	t.Run("should delete tree successfully", func(t *testing.T) {
		app := fiber.New()
		mockTreeService := serviceMock.NewMockTreeService(t)
		app.Delete("/v1/tree/:id", tree.DeleteTree(mockTreeService))
		treeID := int32(1)

		mockTreeService.EXPECT().Delete(
			mock.Anything,
			treeID,
		).Return(nil)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, "/v1/tree/"+strconv.Itoa(int(treeID)), nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		mockTreeService.AssertExpectations(t)
	})

	t.Run("should return 400 Bad Request for invalid tree ID", func(t *testing.T) {
		app := fiber.New()
		mockTreeService := serviceMock.NewMockTreeService(t)
		app.Delete("/v1/tree/:id", tree.DeleteTree(mockTreeService))

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, "/v1/tree/invalid", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 404 Not Found when tree does not exist", func(t *testing.T) {
		app := fiber.New()
		mockTreeService := serviceMock.NewMockTreeService(t)
		app.Delete("/v1/tree/:id", tree.DeleteTree(mockTreeService))

		treeID := int32(999)
		mockTreeService.EXPECT().Delete(
			mock.Anything,
			treeID,
		).Return(service.NewError(service.NotFound, "tree not found"))

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, "/v1/tree/"+strconv.Itoa(int(treeID)), nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockTreeService.AssertExpectations(t)
	})

	t.Run("should return 500 Internal Server Error on service error", func(t *testing.T) {
		app := fiber.New()
		mockTreeService := serviceMock.NewMockTreeService(t)
		app.Delete("/v1/tree/:id", tree.DeleteTree(mockTreeService))

		treeID := int32(1)
		mockTreeService.EXPECT().Delete(
			mock.Anything,
			treeID,
		).Return(fiber.NewError(fiber.StatusInternalServerError, "service error"))

		// when
		body, _ := json.Marshal(TestTreeUpdateRequest)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, "/v1/tree/"+strconv.Itoa(int(treeID)), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockTreeService.AssertExpectations(t)
	})
}

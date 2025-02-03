package wateringplan_test

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
	wateringplan "github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/watering_plan"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	serviceMock "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllWateringPlans(t *testing.T) {
	t.Run("should return all watering plans successfully", func(t *testing.T) {
		app := fiber.New()
		mockWateringPlanService := serviceMock.NewMockWateringPlanService(t)
		handler := wateringplan.GetAllWateringPlans(mockWateringPlanService)
		app.Get("/v1/watering-plan", handler)

		mockWateringPlanService.EXPECT().GetAll(
			mock.Anything,
			"",
		).Return(TestWateringPlans, nil)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/watering-plan", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response serverEntities.WateringPlanListResponse
		err = utils.ParseJSONResponse(resp, &response)
		assert.NoError(t, err)
		assert.Equal(t, len(TestWateringPlans), len(response.Data))
		assert.Equal(t, TestWateringPlans[0].Date, response.Data[0].Date)

		mockWateringPlanService.AssertExpectations(t)
	})

	t.Run("should return an empty list when no watering plans are available", func(t *testing.T) {
		app := fiber.New()
		mockWateringPlanService := serviceMock.NewMockWateringPlanService(t)
		handler := wateringplan.GetAllWateringPlans(mockWateringPlanService)
		app.Get("/v1/watering-plan", handler)

		mockWateringPlanService.EXPECT().GetAll(
			mock.Anything,
			"",
		).Return([]*entities.WateringPlan{}, nil)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/watering-plan", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response serverEntities.WateringPlanListResponse
		err = utils.ParseJSONResponse(resp, &response)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(response.Data))

		mockWateringPlanService.AssertExpectations(t)
	})

	t.Run("should return 500 Internal Server Error when service fails", func(t *testing.T) {
		app := fiber.New()
		mockWateringPlanService := serviceMock.NewMockWateringPlanService(t)
		handler := wateringplan.GetAllWateringPlans(mockWateringPlanService)
		app.Get("/v1/watering-plan", handler)

		mockWateringPlanService.EXPECT().GetAll(
			mock.Anything,
			"",
		).Return(nil, fiber.NewError(fiber.StatusInternalServerError, "service error"))

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/watering-plan", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockWateringPlanService.AssertExpectations(t)
	})
}

func TestGetWateringPlanByID(t *testing.T) {
	t.Run("should return watering plan successfully", func(t *testing.T) {
		app := fiber.New()
		mockWateringPlanService := serviceMock.NewMockWateringPlanService(t)
		handler := wateringplan.GetWateringPlanByID(mockWateringPlanService)
		app.Get("/v1/watering-plan/:id", handler)

		mockWateringPlanService.EXPECT().GetByID(
			mock.Anything,
			int32(1),
		).Return(TestWateringPlans[0], nil)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/watering-plan/1", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response serverEntities.WateringPlanResponse
		err = utils.ParseJSONResponse(resp, &response)
		assert.NoError(t, err)
		assert.Equal(t, TestWateringPlans[0].Date, response.Date)

		mockWateringPlanService.AssertExpectations(t)
	})

	t.Run("should return 400 Bad Request for invalid watering plan ID", func(t *testing.T) {
		app := fiber.New()
		mockWateringPlanService := serviceMock.NewMockWateringPlanService(t)
		handler := wateringplan.GetWateringPlanByID(mockWateringPlanService)
		app.Get("/v1/watering-plan/:id", handler)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/watering-plan/invalid", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 404 Not Found if watering plan does not exist", func(t *testing.T) {
		app := fiber.New()
		mockWateringPlanService := serviceMock.NewMockWateringPlanService(t)
		handler := wateringplan.GetWateringPlanByID(mockWateringPlanService)
		app.Get("/v1/watering-plan/:id", handler)

		mockWateringPlanService.EXPECT().GetByID(
			mock.Anything,
			int32(999),
		).Return(nil, service.NewError(service.NotFound, "not found"))

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/watering-plan/999", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		mockWateringPlanService.AssertExpectations(t)
	})

	t.Run("should return 500 Internal Server Error for service failure", func(t *testing.T) {
		app := fiber.New()
		mockWateringPlanService := serviceMock.NewMockWateringPlanService(t)
		handler := wateringplan.GetWateringPlanByID(mockWateringPlanService)
		app.Get("/v1/watering-plan/:id", handler)

		mockWateringPlanService.EXPECT().GetByID(
			mock.Anything,
			int32(1),
		).Return(nil, fiber.NewError(fiber.StatusInternalServerError, "service error"))

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/watering-plan/1", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockWateringPlanService.AssertExpectations(t)
	})
}

func TestCreateWateringPlan(t *testing.T) {
	t.Run("should create watering plan successfully", func(t *testing.T) {
		app := fiber.New()
		mockWateringPlanService := serviceMock.NewMockWateringPlanService(t)
		handler := wateringplan.CreateWateringPlan(mockWateringPlanService)
		app.Post("/v1/watering-plan", handler)

		mockWateringPlanService.EXPECT().Create(
			mock.Anything,
			mock.AnythingOfType("*entities.WateringPlanCreate"),
		).Return(TestWateringPlans[0], nil)

		// when
		body, _ := json.Marshal(TestWateringPlanRequest)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/v1/watering-plan", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response serverEntities.WateringPlanResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, TestWateringPlans[0].Date, response.Date)
		assert.Equal(t, TestWateringPlans[0].Transporter.NumberPlate, response.Transporter.NumberPlate)

		mockWateringPlanService.AssertExpectations(t)
	})

	t.Run("should return 400 Bad Request for invalid request body", func(t *testing.T) {
		app := fiber.New()
		mockWateringPlanService := serviceMock.NewMockWateringPlanService(t)
		handler := wateringplan.CreateWateringPlan(mockWateringPlanService)
		app.Post("/v1/watering-plan", handler)

		invalidRequestBody := []byte(`{"invalid_field": "value"}`)

		// when
		body, _ := json.Marshal(invalidRequestBody)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/v1/watering-plan", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 500 Internal Server Error for service failure", func(t *testing.T) {
		app := fiber.New()
		mockWateringPlanService := serviceMock.NewMockWateringPlanService(t)
		handler := wateringplan.CreateWateringPlan(mockWateringPlanService)
		app.Post("/v1/watering-plan", handler)

		mockWateringPlanService.EXPECT().Create(
			mock.Anything,
			mock.AnythingOfType("*entities.WateringPlanCreate"),
		).Return(nil, fiber.NewError(fiber.StatusInternalServerError, "service error"))

		// when
		body, _ := json.Marshal(TestWateringPlanRequest)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/v1/watering-plan", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockWateringPlanService.AssertExpectations(t)
	})
}

func TestUpdateWateringPlan(t *testing.T) {
	t.Run("should update watering plan successfully", func(t *testing.T) {
		app := fiber.New()
		mockWateringPlanService := serviceMock.NewMockWateringPlanService(t)
		handler := wateringplan.UpdateWateringPlan(mockWateringPlanService)
		app.Put("/v1/watering-plan/:id", handler)

		mockWateringPlanService.EXPECT().Update(
			mock.Anything,
			int32(1),
			mock.Anything,
		).Return(TestWateringPlans[0], nil)

		// when
		body, _ := json.Marshal(TestWateringPlanRequest)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, "/v1/watering-plan/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response serverEntities.WateringPlanResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, TestWateringPlans[0].Date, response.Date)

		mockWateringPlanService.AssertExpectations(t)
	})

	t.Run("should return 400 Bad Request for invalid watering plan ID", func(t *testing.T) {
		app := fiber.New()
		mockWateringPlanService := serviceMock.NewMockWateringPlanService(t)
		handler := wateringplan.UpdateWateringPlan(mockWateringPlanService)
		app.Put("/v1/watering-plan/:id", handler)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, "/v1/watering-plan/invalid", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 Bad Request for invalid request body", func(t *testing.T) {
		app := fiber.New()
		mockWateringPlanService := serviceMock.NewMockWateringPlanService(t)
		handler := wateringplan.UpdateWateringPlan(mockWateringPlanService)
		app.Put("/v1/watering-plan/:id", handler)

		invalidRequestBody := []byte(`{"invalid_field": "value"}`)

		// when
		body, _ := json.Marshal(invalidRequestBody)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, "/v1/watering-plan/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 404 Not Found if watering plan does not exist", func(t *testing.T) {
		app := fiber.New()
		mockWateringPlanService := serviceMock.NewMockWateringPlanService(t)
		handler := wateringplan.UpdateWateringPlan(mockWateringPlanService)
		app.Put("/v1/watering-plan/:id", handler)

		mockWateringPlanService.EXPECT().Update(
			mock.Anything,
			int32(1),
			mock.Anything,
		).Return(nil, service.NewError(service.NotFound, "not found"))

		// when
		body, _ := json.Marshal(TestWateringPlanRequest)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, "/v1/watering-plan/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		mockWateringPlanService.AssertExpectations(t)
	})

	t.Run("should return 500 Internal Server Error for service failure", func(t *testing.T) {
		app := fiber.New()
		mockWateringPlanService := serviceMock.NewMockWateringPlanService(t)
		handler := wateringplan.UpdateWateringPlan(mockWateringPlanService)
		app.Put("/v1/watering-plan/:id", handler)

		mockWateringPlanService.EXPECT().Update(mock.Anything, int32(1), mock.Anything).Return(nil, fiber.NewError(fiber.StatusInternalServerError, "service error"))

		// when
		body, _ := json.Marshal(TestWateringPlanRequest)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, "/v1/watering-plan/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	})
}

func TestDeleteWateringPlan(t *testing.T) {
	t.Run("should delete watering plan successfully", func(t *testing.T) {
		app := fiber.New()
		mockWateringPlanService := serviceMock.NewMockWateringPlanService(t)
		handler := wateringplan.DeleteWateringPlan(mockWateringPlanService)
		app.Delete("/v1/watering-plan/:id", handler)

		wateringPlanID := 1
		mockWateringPlanService.EXPECT().Delete(mock.Anything, int32(wateringPlanID)).Return(nil)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, "/v1/watering-plan/"+strconv.Itoa(wateringPlanID), nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		mockWateringPlanService.AssertExpectations(t)
	})

	t.Run("should return 400 for invalid ID format", func(t *testing.T) {
		app := fiber.New()
		mockWateringPlanService := serviceMock.NewMockWateringPlanService(t)
		handler := wateringplan.DeleteWateringPlan(mockWateringPlanService)
		app.Delete("/v1/watering-plan/:id", handler)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, "/v1/watering-plan/invalid_id", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 404 for non-existing tree cluster", func(t *testing.T) {
		app := fiber.New()
		mockWateringPlanService := serviceMock.NewMockWateringPlanService(t)
		handler := wateringplan.DeleteWateringPlan(mockWateringPlanService)
		app.Delete("/v1/watering-plan/:id", handler)

		wateringPlanID := 999
		mockWateringPlanService.EXPECT().Delete(
			mock.Anything,
			int32(wateringPlanID),
		).Return(service.NewError(service.NotFound, "not found"))

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, "/v1/watering-plan/"+strconv.Itoa(wateringPlanID), nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockWateringPlanService.AssertExpectations(t)
	})
}

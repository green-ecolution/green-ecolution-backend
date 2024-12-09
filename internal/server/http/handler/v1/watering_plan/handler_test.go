package wateringplan_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	serverEntities "github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	wateringplan "github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/watering_plan"
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
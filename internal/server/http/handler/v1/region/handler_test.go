package region_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	serverEntities "github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/region"
	serviceMock "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllRegions(t *testing.T) {
	t.Run("should return all regions successfully", func(t *testing.T) {
		mockRegionService := serviceMock.NewMockRegionService(t)
		app := fiber.New()
		handler := region.GetAllRegions(mockRegionService)

		expectedRegions := []*entities.Region{
			{ID: 1, Name: "Region A"},
			{ID: 2, Name: "Region B"},
		}

		mockRegionService.On("GetAll", mock.Anything).Return(expectedRegions, nil)
		app.Get("/v1/region", handler)

		// when
		req, _ := http.NewRequest("GET", "/v1/region", nil)
		resp, err := app.Test(req, -1)

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response serverEntities.RegionListResponse
		err = parseJSONResponse(resp, &response)
		assert.NoError(t, err)
		assert.Len(t, response.Regions, 2)
		assert.Equal(t, "Region A", response.Regions[0].Name)
		assert.Equal(t, "Region B", response.Regions[1].Name)

		mockRegionService.AssertExpectations(t)
	})

	t.Run("should return empty region list when no regions found", func(t *testing.T) {
		mockRegionService := serviceMock.NewMockRegionService(t)
		app := fiber.New()
		handler := region.GetAllRegions(mockRegionService)
	
		mockRegionService.On("GetAll", mock.Anything).Return([]*entities.Region{}, nil)
		app.Get("/v1/region", handler)
	
		// when
		req, _ := http.NewRequest("GET", "/v1/region", nil)
		resp, err := app.Test(req, -1)
	
		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	
		var response serverEntities.RegionListResponse
		err = parseJSONResponse(resp, &response)
		assert.NoError(t, err)
		assert.Len(t, response.Regions, 0)
	
		mockRegionService.AssertExpectations(t)
	})

	t.Run("should return 500 when service returns an error", func(t *testing.T) {
		mockRegionService := serviceMock.NewMockRegionService(t)
		app := fiber.New()
		handler := region.GetAllRegions(mockRegionService)
	
		mockRegionService.On("GetAll", mock.Anything).Return(nil, errors.New("service error"))
		app.Get("/v1/region", handler)
	
		// when
		req, _ := http.NewRequest("GET", "/v1/region", nil)
		resp, err := app.Test(req, -1)
	
		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		
		mockRegionService.AssertExpectations(t)
	})
}

func TestGetRegionByID(t *testing.T) {
	t.Run("should return region successfully", func(t *testing.T) {
		mockRegionService := serviceMock.NewMockRegionService(t)
		app := fiber.New()
		handler := region.GetRegionByID(mockRegionService)

		expectedRegion := &entities.Region{
			ID:   1,
			Name: "Region A",
		}

		mockRegionService.On("GetByID", mock.Anything, int32(1)).Return(expectedRegion, nil)
		app.Get("/v1/region/:id", handler)

		// when
		req, _ := http.NewRequest("GET", "/v1/region/1", nil)
		resp, err := app.Test(req, -1)

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response entities.Region
		err = parseJSONResponse(resp, &response)
		assert.NoError(t, err)
		assert.Equal(t, "Region A", response.Name)

		mockRegionService.AssertExpectations(t)
	})

	t.Run("should return 400 for invalid ID", func(t *testing.T) {
		mockRegionService := serviceMock.NewMockRegionService(t)
		app := fiber.New()
		handler := region.GetRegionByID(mockRegionService)

		app.Get("/v1/region/:id", handler)

		// when
		req, _ := http.NewRequest("GET", "/v1/region/invalid-id", nil)
		resp, err := app.Test(req, -1)

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		mockRegionService.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
	})

	t.Run("should return 404 when region not found", func(t *testing.T) {
		mockRegionService := serviceMock.NewMockRegionService(t)
		app := fiber.New()
		handler := region.GetRegionByID(mockRegionService)

		mockRegionService.On("GetByID", mock.Anything, int32(1)).Return(nil, storage.ErrRegionNotFound)
		app.Get("/v1/region/:id", handler)

		// when
		req, _ := http.NewRequest("GET", "/v1/region/1", nil)
		resp, err := app.Test(req, -1)

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		mockRegionService.AssertExpectations(t)
	})

	t.Run("should return 500 when service returns an error", func(t *testing.T) {
		mockRegionService := serviceMock.NewMockRegionService(t)
		app := fiber.New()
		handler := region.GetRegionByID(mockRegionService)

		mockRegionService.On("GetByID", mock.Anything, int32(1)).Return(nil, errors.New("service error"))
		app.Get("/v1/region/:id", handler)

		// when
		req, _ := http.NewRequest("GET", "/v1/region/1", nil)
		resp, err := app.Test(req, -1)

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockRegionService.AssertExpectations(t)
	})
}


// helper function to decode JSON
func parseJSONResponse(body *http.Response, target interface{}) error {
	defer body.Body.Close()
	return json.NewDecoder(body.Body).Decode(target)
}

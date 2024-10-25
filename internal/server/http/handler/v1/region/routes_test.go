package region

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	serverEntities "github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	serviceMock "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterRoutes(t *testing.T) {
	t.Run("GET /v1/region", func(t *testing.T) {
		t.Run("should return all regions successfully", func(t *testing.T) {
			mockRegionService := serviceMock.NewMockRegionService(t)
			app := RegisterRoutes(mockRegionService)

			expectedRegions := []*entities.Region{
				{ID: 1, Name: "Region A"},
				{ID: 2, Name: "Region B"},
			}

			mockRegionService.EXPECT().GetAll(mock.Anything).Return(expectedRegions, nil)

			req, _ := http.NewRequestWithContext(context.Background(), "GET", "/", nil)
			resp, err := app.Test(req, -1)
			defer resp.Body.Close()
			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var response serverEntities.RegionListResponse
			err = utils.ParseJSONResponse(resp, &response)
			assert.NoError(t, err)
			assert.Len(t, response.Regions, 2)
			assert.Equal(t, "Region A", response.Regions[0].Name)
			assert.Equal(t, "Region B", response.Regions[1].Name)

			mockRegionService.AssertExpectations(t)
		})

		t.Run("should return empty list when no regions found", func(t *testing.T) {
			mockRegionService := serviceMock.NewMockRegionService(t)
			app := RegisterRoutes(mockRegionService)

			mockRegionService.EXPECT().GetAll(mock.Anything).Return([]*entities.Region{}, nil)

			req, _ := http.NewRequestWithContext(context.Background(), "GET", "/", nil)
			resp, err := app.Test(req, -1)
			defer resp.Body.Close()
			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var response serverEntities.RegionListResponse
			err = utils.ParseJSONResponse(resp, &response)
			assert.NoError(t, err)
			assert.Len(t, response.Regions, 0)

			mockRegionService.AssertExpectations(t)
		})
	})

	t.Run("GET /v1/region/:id", func(t *testing.T) {
		t.Run("should return region by ID successfully", func(t *testing.T) {
			mockRegionService := serviceMock.NewMockRegionService(t)
			app := RegisterRoutes(mockRegionService)

			expectedRegion := &entities.Region{ID: 1, Name: "Region A"}

			mockRegionService.EXPECT().GetByID(mock.Anything, int32(1)).Return(expectedRegion, nil)

			req, _ := http.NewRequestWithContext(context.Background(), "GET", "/1", nil)
			resp, err := app.Test(req, -1)
			defer resp.Body.Close()
			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var response entities.Region
			err = utils.ParseJSONResponse(resp, &response)
			assert.NoError(t, err)
			assert.Equal(t, "Region A", response.Name)

			mockRegionService.AssertExpectations(t)
		})

		t.Run("should return 400 for invalid ID", func(t *testing.T) {
			mockRegionService := serviceMock.NewMockRegionService(t)
			app := RegisterRoutes(mockRegionService)

			req, _ := http.NewRequestWithContext(context.Background(), "GET", "/invalid-id", nil)
			resp, err := app.Test(req, -1)
			defer resp.Body.Close()
			assert.Nil(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

			mockRegionService.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
		})

		t.Run("should return 404 when region not found", func(t *testing.T) {
			mockRegionService := serviceMock.NewMockRegionService(t)
			app := RegisterRoutes(mockRegionService)

			mockRegionService.EXPECT().GetByID(mock.Anything, int32(1)).Return(nil, storage.ErrRegionNotFound)

			req, _ := http.NewRequestWithContext(context.Background(), "GET", "/1", nil)
			resp, err := app.Test(req, -1)
			defer resp.Body.Close()
			assert.Nil(t, err)
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)

			mockRegionService.AssertExpectations(t)
		})

		t.Run("should return 500 when service error occurs", func(t *testing.T) {
			mockRegionService := serviceMock.NewMockRegionService(t)
			app := RegisterRoutes(mockRegionService)

			mockRegionService.EXPECT().GetByID(mock.Anything, int32(1)).Return(nil, errors.New("service error"))

			req, _ := http.NewRequestWithContext(context.Background(), "GET", "/1", nil)
			resp, err := app.Test(req, -1)
			defer resp.Body.Close()
			assert.Nil(t, err)
			assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

			mockRegionService.AssertExpectations(t)
		})
	})
}

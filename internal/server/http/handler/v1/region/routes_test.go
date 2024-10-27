package region

import (
	"net/http"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	serviceMock "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterRoutes(t *testing.T) {
	t.Run("/v1/region/", func(t *testing.T) {
		t.Run("should call GET handler", func(t *testing.T) {
			mockRegionService := serviceMock.NewMockRegionService(t)
			app := RegisterRoutes(mockRegionService)

			// when
			req, _ := http.NewRequest(http.MethodGet, "/", nil)

			expectedRegions := []*entities.Region{
				{ID: 1, Name: "Region A"},
				{ID: 2, Name: "Region B"},
			}

			mockRegionService.EXPECT().GetAll(mock.Anything).Return(expectedRegions, nil)
			
			// then
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	})

	t.Run("/v1/region/:id", func(t *testing.T) {
		t.Run("should call GET handler", func(t *testing.T) {
			mockRegionService := serviceMock.NewMockRegionService(t)
			app := RegisterRoutes(mockRegionService)
			
			// when
			req, _ := http.NewRequest(http.MethodGet, "/1", nil)

			expectedRegion := &entities.Region{ID: 1, Name: "Region A"}
			
			mockRegionService.EXPECT().GetByID(mock.Anything, int32(1)).Return(expectedRegion, nil)
			
			// then
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	})
}
